package match

import (
	"bytes"
	"compress/gzip"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
	"match/src/constants"
	"match/src/domain"
	"match/src/log"
	"match/src/pb"
	"match/src/redis"
	"math/rand"
	"strconv"
	"time"
)

const (
	GAME_FPS            = 8
	GAME_FRAME_INTERVAL = time.Second / GAME_FPS
)

type GameMatch struct {
	Server      *Server
	log         *log.Logger
	Config      *GlobalConfig
	gameId      string
	RedisDao    *redis.RedisDao
	redisKey    string
	lockKey     string
	roomKey     string
	GameInfo    *domain.GameInfo
	frameTicker *time.Ticker
	FrameID     int64
	NextFrameId int64
	rand        *rand.Rand
	Counter     *AtomicCounter
	MsgFromNats chan Closure
	exit        chan bool
}

func NewGameMatch(server *Server, gameInfo *domain.GameInfo) *GameMatch {
	now := time.Now()
	return &GameMatch{
		Server:      server,
		log:         server.Log,
		Config:      server.Config,
		gameId:      gameInfo.GameId,
		GameInfo:    gameInfo,
		RedisDao:    server.RedisDao,
		redisKey:    gameInfo.GameId + ".match.queue",
		lockKey:     gameInfo.GameId + ".match.lock",
		roomKey:     "match_room_key",
		FrameID:     0,
		NextFrameId: 0,
		frameTicker: nil,
		MsgFromNats: make(chan Closure, 10*1024),
		rand:        rand.New(rand.NewSource(now.Unix())),
		Counter:     &AtomicCounter{},
		exit:        make(chan bool, 1),
	}
}

func (self *GameMatch) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	self.log.Info("game match %+v begin ...", self.gameId)

	self.frameTicker = time.NewTicker(GAME_FRAME_INTERVAL)
	defer func() {
		self.frameTicker.Stop()
	}()

	self.Server.WaitGroup.Add(1)
	defer func() {
		self.Server.WaitGroup.Done()
	}()

ALL:
	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			return
		case c, ok := <-self.MsgFromNats:
			if !ok {
				break ALL
			}
			SafeRunClosure(self, c)
		case <-self.frameTicker.C:
			SafeRunClosure(self, func() {
				self.Frame()
			})
		default:
			// do nothing
		}
	}
}

func (self *GameMatch) Frame() {
	self.FrameID++

	if self.FrameID < self.NextFrameId {
		return
	} else {
		self.NextFrameId = self.FrameID + GAME_FPS + RandomInt(self.rand, 0, 10)
	}

	lock, err := self.RedisDao.Lock(self.lockKey, "1", 2)
	if err != nil {
		self.log.Error("game match redis lock error: %+v", err)
		return
	}
	if lock {
		for {
			count, err := self.RedisDao.LLen(self.redisKey)
			if err != nil {
				self.RedisDao.Unlock(self.lockKey)
				self.log.Error("redis count err: %v", err)
				return
			}
			if count == 0 {
				self.RedisDao.Unlock(self.lockKey)
				return
			}

			if count >= 2 {
				firstData, _ := self.RedisDao.RPop(self.redisKey)
				secondData, _ := self.RedisDao.RPop(self.redisKey)

				firstMatchReq := self.UnCompressed([]byte(firstData))
				secondMatchReq := self.UnCompressed([]byte(secondData))
				uidList := make([]*pb.MatchData, 0)
				data := &pb.MatchData{
					Pid: firstMatchReq.GetPid(),
					Uid: firstMatchReq.GetUid(),
				}
				uidList = append(uidList, data)

				secondUidData := &pb.MatchData{
					Pid: secondMatchReq.GetPid(),
					Uid: secondMatchReq.GetUid(),
				}
				uidList = append(uidList, secondUidData)

				roomId := self.GenRoomId()
				gameIp := self.PublicCreateRoom(roomId)
				self.Send2PlayerMatchRes(firstMatchReq, uidList, roomId, gameIp)
				self.Send2PlayerMatchRes(secondMatchReq, uidList, roomId, gameIp)
			} else {
				//判断匹配时间
				data, _ := self.RedisDao.RPop(self.redisKey)
				matchReq := self.UnCompressed([]byte(data))
				if matchReq.GetTimeStamp()+int64(self.GameInfo.MatchTime) > int64(time.Now().Unix()) {
					roomId := self.GenRoomId()
					gameIp := self.PublicCreateRoom(roomId)
					self.Send2PlayerWithAI(matchReq, roomId, gameIp)
				} else {
					self.RedisDao.LPush(self.redisKey, data)
				}

				break
			}
		}

		self.RedisDao.Unlock(self.lockKey)
	}
}

func (self *GameMatch) PublicCreateRoom(roomId string) string {
	subject := constants.GetCreateRoomNoticeSubject(self.GameInfo.GameId, self.GameInfo.GroupName)
	var response interface{}
	self.log.Info("PublicCreateRoom %+v roomId %+v", subject, roomId)
	self.Server.NatsPool.Request(subject, roomId, &response, 1*time.Second)
	self.log.Info("PublicCreateRoom roomId %+v response: %+v", roomId, response)
	dataMap := response.(map[string]interface{})
	hostIp := dataMap["data"].(string)

	return hostIp
}

func (self *GameMatch) GenRoomId() string {
	roomId, _ := self.RedisDao.IncrBy(self.roomKey, 1)
	return self.gameId + "_room_" + strconv.FormatInt(roomId+1000, 10)
}

func (self *GameMatch) Send2PlayerWithAI(matchReq *pb.MatchRequest, roomId string, gameIp string) {
	uidList := make([]*pb.MatchData, 0)
	matchData := &pb.MatchData{
		Pid: matchReq.Pid,
		Uid: matchReq.Uid,
	}
	uidList = append(uidList, matchData)
	aiList := make([]*pb.MatchData, 0)
	aiInfo := self.GetAiInfo()
	if aiInfo == nil {
		return
	}

	aiData := &pb.MatchData{
		Pid: aiInfo.Pid,
		Uid: aiInfo.Uid,
	}

	aiList = append(aiList, aiData)
	response := &pb.MatchOverRes{
		Code:      constants.CODE_SUCCESS,
		Msg:       "Success",
		GameId:    self.gameId,
		Uid:       matchReq.GetUid(),
		RoomId:    roomId,
		Timestamp: time.Now().Unix(),
		UidList:   uidList,
		AiUidList: aiList,
		GameIp:    gameIp,
	}

	res, _ := proto.Marshal(response)

	self.Server.NatsPool.Publish(matchReq.GetReceiveSubject(), string(res))
}

func (self *GameMatch) Send2PlayerMatchRes(matchReq *pb.MatchRequest, uidList []*pb.MatchData, roomId string, gameIp string) {
	response := &pb.MatchOverRes{
		Code:      constants.CODE_SUCCESS,
		Msg:       "Success",
		GameId:    self.gameId,
		Uid:       matchReq.GetUid(),
		RoomId:    roomId,
		Timestamp: time.Now().Unix(),
		UidList:   uidList,
		AiUidList: make([]*pb.MatchData, 0),
		GameIp:    gameIp,
	}

	res, _ := proto.Marshal(response)
	self.Server.NatsPool.Publish(matchReq.GetReceiveSubject(), string(res))
}

func (self *GameMatch) Compressed(req *pb.MatchRequest) []byte {
	// 序列化 protobuf 消息
	data, err := proto.Marshal(req)
	if err != nil {
		self.log.Error("Failed to marshal message: %+v", err)
	}
	var compressedData []byte
	buf := bytes.NewBuffer(compressedData)
	writer := gzip.NewWriter(buf)
	_, err = writer.Write(data)
	if err != nil {
		self.log.Error("Failed to compress data: %+v", err)
	}
	writer.Close()
	return buf.Bytes()
}

func (self *GameMatch) UnCompressed(compressedData []byte) *pb.MatchRequest {
	buf := bytes.NewBuffer(compressedData)
	reader, err := gzip.NewReader(buf)
	if err != nil {
		self.log.Error("Failed to create gzip reader: %+v", err)
	}
	data, err := ioutil.ReadAll(reader)
	if err != nil {
		self.log.Error("Failed to read compressed data: %v", err)
	}

	// 反序列化 protobuf 消息
	var msg pb.MatchRequest
	err = proto.Unmarshal(data, &msg)
	if err != nil {
		self.log.Error("Failed to unmarshal message: %v", err)
	}

	reader.Close()
	return &msg
}

func (self *GameMatch) AddMatchRequest(req *pb.MatchRequest) {
	//单机游戏， 直接生成一个roomId，返回
	uidList := make([]*pb.MatchData, 0)
	data := &pb.MatchData{
		Pid: req.GetPid(),
		Uid: req.GetUid(),
	}
	uidList = append(uidList, data)
	gameType := self.GameInfo.Type
	if gameType == constants.GAME_TYPE_SINGLE || gameType == constants.GAME_TYPE_HALL_SINGLE {
		response := &pb.MatchOverRes{
			Code:      constants.CODE_SUCCESS,
			Msg:       "Success",
			GameId:    self.gameId,
			Uid:       req.GetUid(),
			RoomId:    self.gameId + "_room_" + req.GetUid(),
			Timestamp: time.Now().Unix(),
			UidList:   uidList,
			AiUidList: make([]*pb.MatchData, 0),
		}

		res, _ := proto.Marshal(response)
		self.Server.NatsPool.Publish(req.GetReceiveSubject(), string(res))
	} else {
		buf := self.Compressed(req)
		// 将压缩后的数据存储到 Redis 中
		self.RedisDao.LPush(self.redisKey, buf)
	}
}

func (self *GameMatch) CancelMatchRequest(req *pb.CancelMatchRequest) {
	lock, err := self.RedisDao.Lock(self.lockKey, "1", 3)
	if err != nil {
		self.log.Error("CancelMatchRequest game match redis lock error: %+v", err)
		return
	}
	if lock {
		count, err := self.RedisDao.LLen(self.redisKey)
		if err != nil {
			self.RedisDao.Unlock(self.lockKey)
			self.log.Error("redis count err: %v", err)
			return
		}
		if count == 0 {
			self.RedisDao.Unlock(self.lockKey)
			return
		}

		matchList := make([]*pb.MatchRequest, 0)
		for {
			data, err := self.RedisDao.RPop(self.redisKey)
			if err != nil || len(data) < 1 {
				break
			}
			if len(data) > 1 {
				matchReq := self.UnCompressed([]byte(data))
				if matchReq.GetUid() == req.GetUid() {
					break
				}
				matchList = append(matchList, matchReq)
			}
		}

		for _, matchReq := range matchList {
			buf := self.Compressed(matchReq)
			self.RedisDao.LPush(self.redisKey, buf)
		}
	}
}

func (self *GameMatch) LoginHallRequest(reply string, req *pb.CreateHallRequest) {
	self.log.Info("LoginHallRequest  req %+v game %+v", req.GetUid(), req.GetGameId())
	roomId := self.gameId + "_hall_" + req.GetUid()

	gameIp := self.PublicCreateRoom(roomId)

	response := &pb.CreateHallResponse{
		Code:   constants.CODE_SUCCESS,
		Msg:    "Success",
		GameId: self.gameId,
		Uid:    req.GetUid(),
		RoomId: roomId,
		GameIp: gameIp,
	}

	res, _ := proto.Marshal(response)
	self.Server.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
}
func (self *GameMatch) Quit() {
	self.exit <- true
}

/*
*
获取用户所在的分区， 如果是染色用户， 就根据染色设置，否则就是 game-info的分区
*/
func (self *GameMatch) GetGroupName(uid string) string {
	//判断是否是染色用户
	groupName := self.GameInfo.GroupName
	setKey := "COLORED_UID_SET_KEY" + self.GameInfo.GameId
	isMember, err := self.RedisDao.SISMMBER(setKey, uid)
	if err != nil {
		return groupName
	}

	if isMember {
		redisKey := "COLORED_UID_KEY" + self.GameInfo.GameId
		data, err := self.RedisDao.Get(redisKey)
		if err != nil {
			return groupName
		}

		groupName = data
	}

	return groupName
}

func (self *GameMatch) GetAiInfo() *domain.AiInfo {
	aiSetKey := "AI_UID_PID_SET_KEY"
	aiHKey := "AI_UID_PID_H_KEY"
	uid, err := self.RedisDao.SRandMember(aiSetKey)
	if err != nil || uid == "" {
		self.log.Error("GetAiInfo err: %v", err)
		return nil
	}

	pid, err := self.RedisDao.HGet(aiHKey, uid)
	if err != nil {
		self.log.Error("GetAiInfo err: %v", err)
		return nil
	}

	res := &domain.AiInfo{
		Uid: uid,
		Pid: pid,
	}

	return res
}
