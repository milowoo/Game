package match

import (
	"github.com/golang/protobuf/proto"
	"match/src/config"
	"match/src/constants"
	"match/src/domain"
	"match/src/internal"
	"match/src/pb"
	"match/src/utils"
	"math/rand"
	"strconv"
	"time"
)

const (
	GAME_FPS            = 8
	GAME_FRAME_INTERVAL = time.Second / GAME_FPS
	GAME_MATCH_TIME     = time.Second * 2
)

type GameMatch struct {
	Server       *Server
	Config       *config.GlobalConfig
	gameId       string
	redisKey     string
	lockKey      string
	lastMatchKey string
	roomKey      string
	GameInfo     *domain.GameInfo
	frameTicker  *time.Ticker
	FrameID      int64
	NextFrameId  int64
	rand         *rand.Rand
	Counter      *utils.AtomicCounter
	MsgFromNats  chan utils.Closure
	exit         chan bool
}

func NewGameMatch(server *Server, gameInfo *domain.GameInfo) *GameMatch {
	now := time.Now()
	return &GameMatch{
		Server:       server,
		Config:       server.Config,
		gameId:       gameInfo.GameId,
		GameInfo:     gameInfo,
		redisKey:     gameInfo.GameId + ".match.queue",
		lockKey:      gameInfo.GameId + ".match.lock",
		lastMatchKey: gameInfo.GameId + ".match.last.time",
		roomKey:      "match_room_key",
		FrameID:      0,
		NextFrameId:  0,
		frameTicker:  nil,
		MsgFromNats:  make(chan utils.Closure, 10*1024),
		rand:         rand.New(rand.NewSource(now.Unix())),
		Counter:      &utils.AtomicCounter{},
		exit:         make(chan bool, 1),
	}
}

func (self *GameMatch) Run() {
	defer func() {
		p := recover()
		if p != nil {
			internal.GLog.Info("execute panic recovered and going to stop: %v", p)
		}
	}()

	internal.GLog.Info("game match %+v begin ...", self.gameId)

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
			utils.SafeRunClosure(self, c)
		case <-self.frameTicker.C:
			utils.SafeRunClosure(self, func() {
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
	}

	self.NextFrameId = self.FrameID + int64(GAME_MATCH_TIME/GAME_FRAME_INTERVAL)
	lock, err := internal.RedisDao.Lock(self.lockKey, "1", 2*time.Second)
	if err != nil {
		internal.GLog.Error("game match redis lock error: %+v", err)
		return
	}
	if lock {
		lastMatchTime, _ := internal.RedisDao.Get(self.lastMatchKey)
		if len(lastMatchTime) > 1 {
			matchTime, _ := strconv.ParseInt(lastMatchTime, 0, 64)
			if time.Now().UnixMilli()-matchTime < 1000 {
				internal.GLog.Info(" last match time is more than 1000 matchTime [%+v] cur time [%+v]",
					matchTime, time.Now().UnixMilli())
				internal.RedisDao.Unlock(self.lockKey)
				return
			}
		}
		for {
			count, _ := internal.RedisDao.LLen(self.redisKey)
			if count == 0 {
				internal.RedisDao.Unlock(self.lockKey)
				return
			}

			if count >= 2 {
				firstData, _ := internal.RedisDao.RPop(self.redisKey)
				secondData, _ := internal.RedisDao.RPop(self.redisKey)
				var firstMatchReq pb.MatchRequest
				var secondMatchReq pb.MatchRequest
				proto.Unmarshal([]byte(firstData), &firstMatchReq)
				proto.Unmarshal([]byte(secondData), &secondMatchReq)

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
				self.Send2PlayerMatchRes(&firstMatchReq, uidList, roomId, gameIp)
				self.Send2PlayerMatchRes(&secondMatchReq, uidList, roomId, gameIp)
				self.deleteMatchUser(self.gameId, firstMatchReq.GetUid())
				self.deleteMatchUser(self.gameId, secondMatchReq.GetUid())
			} else {
				//判断匹配时间
				data, _ := internal.RedisDao.RPop(self.redisKey)
				var matchReq pb.MatchRequest
				proto.Unmarshal([]byte(data), &matchReq)

				internal.GLog.Info("match  cur time [%+v] data %+v ", time.Now().UnixMilli(), &matchReq)
				timestamp, _ := strconv.ParseInt(matchReq.GetTimeStamp(), 0, 64)
				if timestamp+int64(self.GameInfo.MatchTime)*1000 > time.Now().UnixMilli() {
					roomId := self.GenRoomId()
					gameIp := self.PublicCreateRoom(roomId)
					self.Send2PlayerWithAI(&matchReq, roomId, gameIp)
					self.deleteMatchUser(self.gameId, matchReq.GetUid())
				} else {
					internal.RedisDao.LPush(self.redisKey, data)
				}
				break
			}
		}
		curTime := strconv.FormatInt(time.Now().UnixMilli(), 10)
		internal.RedisDao.Set(self.lastMatchKey, curTime, 1000*time.Second)
		internal.RedisDao.Unlock(self.lockKey)
	}
}

func (self *GameMatch) PublicCreateRoom(roomId string) string {
	subject := constants.GetCreateRoomNoticeSubject(self.GameInfo.GameId, self.GameInfo.GroupName)
	var response interface{}
	internal.GLog.Info("PublicCreateRoom %+v roomId %+v", subject, roomId)
	internal.NatsPool.Request(subject, roomId, &response, 1*time.Second)
	if response == nil {
		return ""
	}
	internal.GLog.Info("PublicCreateRoom roomId %+v response: %+v", roomId, response)
	dataMap := response.(map[string]interface{})
	hostIp := dataMap["data"].(string)

	return hostIp
}

func (self *GameMatch) GenRoomId() string {
	roomId, _ := internal.RedisDao.IncrBy(self.roomKey, 1)
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
		internal.GLog.Error("Send2PlayerWithAI get aiInfo failed %")
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
		Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		UidList:   uidList,
		AiUidList: aiList,
		GameIp:    gameIp,
	}

	res, _ := proto.Marshal(response)
	internal.GLog.Info("Send2PlayerWithAI uid %+v ", uidList)
	internal.NatsPool.Publish(matchReq.GetReceiveSubject(), string(res))
}

func (self *GameMatch) Send2PlayerMatchRes(matchReq *pb.MatchRequest, uidList []*pb.MatchData, roomId string, gameIp string) {
	response := &pb.MatchOverRes{
		Code:      constants.CODE_SUCCESS,
		Msg:       "Success",
		GameId:    self.gameId,
		Uid:       matchReq.GetUid(),
		RoomId:    roomId,
		Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		UidList:   uidList,
		AiUidList: make([]*pb.MatchData, 0),
		GameIp:    gameIp,
	}

	internal.GLog.Info("Send2PlayerMatchRes gameId %+v uid %+v", self.gameId, uidList)

	res, _ := proto.Marshal(response)
	internal.NatsPool.Publish(matchReq.GetReceiveSubject(), string(res))
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
		roomId := self.gameId + "_room_" + req.GetUid()
		hostIp := self.PublicCreateRoom(roomId)
		//通知
		response := &pb.MatchOverRes{
			Code:      constants.CODE_SUCCESS,
			Msg:       "Success",
			GameId:    self.gameId,
			Uid:       req.GetUid(),
			RoomId:    roomId,
			Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
			UidList:   uidList,
			AiUidList: make([]*pb.MatchData, 0),
			GameIp:    hostIp,
		}

		res, _ := proto.Marshal(response)
		internal.GLog.Info("AddMatchRequest response subject %+v", req.GetReceiveSubject())
		internal.NatsPool.Publish(req.GetReceiveSubject(), string(res))
		self.deleteMatchUser(self.gameId, req.GetUid())
	} else {
		internal.GLog.Info("AddMatchRequest 222")
		buf, _ := proto.Marshal(req)
		// 将数据存储到 Redis 中
		internal.RedisDao.LPush(self.redisKey, string(buf))
	}

}

func (self *GameMatch) deleteMatchUser(gameId, uid string) {
	redisKey := utils.GetMatchUserKey(gameId, uid)
	internal.RedisDao.Del(redisKey)
}

func (self *GameMatch) CancelMatchRequest(req *pb.CancelMatchRequest) {
	lock, err := internal.RedisDao.Lock(self.lockKey, "1", 3*time.Second)
	if err != nil {
		internal.GLog.Error("CancelMatchRequest game match redis lock error: %+v", err)
		return
	}
	if lock {
		count, err := internal.RedisDao.LLen(self.redisKey)
		if err != nil {
			internal.RedisDao.Unlock(self.lockKey)
			internal.GLog.Error("redis count err: %v", err)
			return
		}
		if count == 0 {
			internal.RedisDao.Unlock(self.lockKey)
			return
		}

		matchList := make([]*pb.MatchRequest, 0)
		for {
			data, err := internal.RedisDao.RPop(self.redisKey)
			if err != nil || len(data) < 1 {
				break
			}
			if len(data) > 1 {
				var matchReq pb.MatchRequest
				proto.Unmarshal([]byte(data), &matchReq)
				if matchReq.GetUid() == req.GetUid() {
					break
				}
				matchList = append(matchList, &matchReq)
			}
		}

		for _, matchReq := range matchList {
			buf, _ := proto.Marshal(matchReq)
			internal.RedisDao.LPush(self.redisKey, string(buf))
		}
	}
}

func (self *GameMatch) LoginHallRequest(reply string, req *pb.CreateHallRequest) {
	internal.GLog.Info("LoginHallRequest  req %+v game %+v", req.GetUid(), req.GetGameId())
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
	internal.NatsPool.Publish(reply, map[string]interface{}{"res": "ok", "data": string(res)})
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
	isMember, err := internal.RedisDao.SISMMBER(setKey, uid)
	if err != nil {
		return groupName
	}

	if isMember {
		redisKey := "COLORED_UID_KEY" + self.GameInfo.GameId
		data, err := internal.RedisDao.Get(redisKey)
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
	uid, err := internal.RedisDao.SRandMember(aiSetKey)
	if err != nil || uid == "" {
		internal.GLog.Error("GetAiInfo err: %v", err)
		return nil
	}

	pid, err := internal.RedisDao.HGet(aiHKey, uid)
	if err != nil {
		internal.GLog.Error("GetAiInfo err: %v", err)
		return nil
	}

	res := &domain.AiInfo{
		Uid: uid,
		Pid: pid,
	}

	return res
}

func RunOnMatch(c chan utils.Closure, mgr *GameMatch, cb func(mgr *GameMatch)) {
	c <- func() {
		cb(mgr)
	}
}
