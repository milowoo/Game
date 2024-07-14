package match

import (
	"bytes"
	"compress/gzip"
	"github.com/golang/protobuf/proto"
	"io/ioutil"
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
	Server            *Server
	log               *log.Logger
	Config            *GlobalConfig
	gameId            string
	gameType          GameType
	RedisDao          *redis.RedisDao
	redisKey          string
	lockKey           string
	roomKey           string
	createRoomSubject string
	GameInfo          *GameInfo
	frameTicker       *time.Ticker
	FrameID           int64
	NextFrameId       int64
	rand              *rand.Rand
	MsgFromNats       chan Closure
	exit              chan bool
}

func NewGameMatch(server *Server, gameInfo *GameInfo) *GameMatch {
	now := time.Now()
	return &GameMatch{
		Server:            server,
		log:               server.Log,
		Config:            server.Config,
		gameId:            gameInfo.GameId,
		GameInfo:          gameInfo,
		gameType:          GameType(gameInfo.Type),
		RedisDao:          server.RedisDao,
		redisKey:          gameInfo.GameId + ".match.queue",
		lockKey:           gameInfo.GameId + ".match.lock",
		roomKey:           "match_room_key",
		createRoomSubject: "create.room.notice." + gameInfo.GameId,
		FrameID:           0,
		NextFrameId:       0,
		frameTicker:       nil,
		MsgFromNats:       make(chan Closure, 2048),
		rand:              rand.New(rand.NewSource(now.Unix())),
		exit:              make(chan bool, 1),
	}
}

func (self *GameMatch) Run() {
	defer func() {
		p := recover()
		if p != nil {
			self.log.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	self.frameTicker = time.NewTicker(GAME_FRAME_INTERVAL)
	defer func() {
		self.frameTicker.Stop()
	}()

	for {
		// 优先查看exit，
		select {
		case <-self.exit:
			return
		default:
			// do nothing
		}

		select {
		case <-self.frameTicker.C:
			SafeRunClosure(self, func() {
				self.Frame()
			})
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
				uidList := make([]string, 0)
				uidList = append(uidList, firstMatchReq.GetUid())
				uidList = append(uidList, secondMatchReq.GetUid())
				roomId := self.GenRoomId()
				self.Send2PlayerMatchRes(firstMatchReq, uidList, roomId)
				self.Send2PlayerMatchRes(secondMatchReq, uidList, roomId)
			} else {
				//判断匹配时间
				data, _ := self.RedisDao.RPop(self.redisKey)
				matchReq := self.UnCompressed([]byte(data))
				if matchReq.GetTimeStamp()+int64(self.GameInfo.MatchTime) > int64(time.Now().Unix()) {
					roomId := self.GenRoomId()
					self.Send2PlayerWithAI(matchReq, roomId)
				} else {
					self.RedisDao.LPush(self.redisKey, data)
				}

				break
			}
		}

		self.RedisDao.Unlock(self.lockKey)
	}
}

func (self *GameMatch) PublicCreateRoom(roomId string) {
	self.Server.NatsPool.Publish(self.createRoomSubject, []byte(roomId))
}

func (self *GameMatch) GenRoomId() string {
	roomId, _ := self.RedisDao.IncrBy(self.roomKey, 1)
	return self.gameId + "_room_" + strconv.FormatInt(roomId+1000, 10)
}

func (self *GameMatch) Send2PlayerWithAI(matchReq *pb.MatchRequest, roomId string) {
	uidList := make([]string, 0)
	uidList = append(uidList, matchReq.GetUid())
	aiList := make([]string, 0)
	aiList = append(aiList, "1001")
	response := &pb.MatchOverRes{
		Code:      100,
		Msg:       "Success",
		GameId:    self.gameId,
		Uid:       matchReq.GetUid(),
		RoomId:    roomId,
		Timestamp: time.Now().Unix(),
		UidList:   uidList,
		AiUidList: aiList,
	}

	res, _ := proto.Marshal(response)
	self.Server.NatsPool.Publish(matchReq.GetReceiveSubject(), res)
}

func (self *GameMatch) Send2PlayerMatchRes(matchReq *pb.MatchRequest, uidList []string, roomId string) {
	response := &pb.MatchOverRes{
		Code:      100,
		Msg:       "Success",
		GameId:    self.gameId,
		Uid:       matchReq.GetUid(),
		RoomId:    roomId,
		Timestamp: time.Now().Unix(),
		UidList:   uidList,
		AiUidList: make([]string, 0),
	}

	res, _ := proto.Marshal(response)
	self.Server.NatsPool.Publish(matchReq.GetReceiveSubject(), res)
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
	uidList := make([]string, 0)
	uidList = append(uidList, req.GetUid())

	if self.gameType == Single || self.gameType == HallAndSingle || (self.gameType == HallAnd1V1 && req.GetOpt() == "HALL") {
		response := &pb.MatchOverRes{
			Code:      100,
			Msg:       "Success",
			GameId:    self.gameId,
			Uid:       req.GetUid(),
			RoomId:    self.gameId + "_room_" + req.GetUid(),
			Timestamp: time.Now().Unix(),
			UidList:   uidList,
			AiUidList: make([]string, 0),
		}

		res, _ := proto.Marshal(response)
		self.Server.NatsPool.Publish(req.GetReceiveSubject(), res)
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

func (self *GameMatch) Quit() {
	self.exit <- true
}
