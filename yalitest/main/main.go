package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/golang/protobuf/proto"

	"log"
	stresstest "yalitest"
	"yalitest/constants"
	log2 "yalitest/log"
	"yalitest/pb"
	stat "yalitest/stat"
	util "yalitest/util"
	websocket "yalitest/websocket"
)

// 消息发送计数器
var messageMatch = stat.NewAsyncRequestCounterDefault()
var messageEnterHall = stat.NewAsyncRequestCounterDefault()
var messageEnterRoom = stat.NewAsyncRequestCounterDefault()

var g_gameId = "test4"
var gLog *log2.Logger
var gCounter *util.AtomicCounter

var baseUid = int64(30000000)

func RandUid() int64 {
	timer := time.NewTimer(0)
	timer.Stop()
	now := time.Now()
	rand1 := rand.New(rand.NewSource(now.Unix()))
	uid := util.RandomInt(rand1, 0, 100000000)

	return int64(uid)
}

func RandString() string {
	uuid := util.UUID()
	for i := 0; i < 20; i++ {
		uuid = uuid + util.UUID()
	}
	return uuid
}

func main() {
	logName := "/Users/wuchuangeng/game/logs/" + "test.log"
	//日志名 + 文件大小（M为单位） + 打印标志 + 线程数量 （未启动） + 工作协程长度（未启动） + 深度
	Log := log2.NewLogger2(logName, 1024*2, log.LstdFlags|log.Lshortfile, 8, 1024, 2)
	gLog = Log

	gCounter = &util.AtomicCounter{}

	gLog.Info("yalitest begin .........")

	reporter := stat.NewReporter()
	reporter.AddCollector("match", messageMatch)
	reporter.AddCollector("enterHall", messageEnterHall)
	reporter.AddCollector("enterRoom", messageEnterRoom)

	clientNum, _ := strconv.ParseInt("0", 10, 64)

	if len(os.Args) >= 2 {
		clientNum, _ = strconv.ParseInt(os.Args[1], 10, 64)
	}

	mgr := stresstest.Manager{
		ClientCount:      clientNum,
		ConnectPerSecond: 40,
		ExecutePerSecond: 1,
		ClientSpawnFunc:  spawnClient,
		Reporter:         reporter,
	}

	err := mgr.Run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "stress manager run fail: %v", err)
	}
}

func spawnClient(n int64) stresstest.Client {
	IntUid := baseUid + n
	//uid := strconv.FormatInt(RandUid(), 10)
	uid := strconv.FormatInt(IntUid, 10)
	u := &websocket.URL{
		Addr: "ws://127.0.0.1:36001",
		//Addr:"wss://xyx-dk-stress-cn-804-2.kaixindou.net",
		GameName: g_gameId,
		AppKey:   "xAZ9JCAYeDlJiS8VB7Q4MRl5G76jcL3f",
		UID:      uid,
	}

	return &MessageClient{
		Uid:            uid,
		BaseSeq:        util.UUID(),
		URL:            u.String(),
		IsEnterRoom:    false,
		IsLoadProgress: false,
		ExcuteNum:      1,
		IsHall:         true,
		IsMatch:        false,
		IsFirst:        true,
		IsGameOver:     false,
	}
}

// MessageClient ...
type MessageClient struct {
	URL            string
	Conn           *websocket.Connection
	Uid            string
	BaseSeq        string
	RoomId         string
	ExcuteNum      int
	IsEnterRoom    bool
	IsHall         bool
	IsLoadProgress bool
	IsMatch        bool
	IsFirst        bool
	IsGameOver     bool
	Log            *log2.Logger
}

func (c *MessageClient) getReportSeq(seq int64) interface{} {
	return fmt.Sprintf("%v_%v", c.BaseSeq, seq)
}

func (c *MessageClient) Receive(down *pb.ClientCommonResponse) {
	//ylog.Info(fmt.Sprintf("uid %v Receive  body Name %v", c.Uid, down.BodyName))
	bodyName := down.GetHead().ProtoName
	if down.Code != 100 {
		gLog.Error("Receive err code %+v msg %+v proto name %+v",
			down.Code, down.Msg, bodyName)
		return
	}

	gLog.Info("Receive body Name %+v", bodyName)

	if bodyName == "pb.ClientLoginHallResponse" {
		var res pb.ClientLoginHallResponse
		proto.Unmarshal(down.Body, &res)
		if res.Code != constants.CODE_SUCCESS {
			gLog.Info("login err %+v msg %+v", res.Code, res.Msg)
			messageEnterHall.AddFail(c.getReportSeq(down.Head.Sn))
			return
		}
		messageEnterHall.AddSuccess(c.getReportSeq(down.Head.Sn))

		c.IsHall = true

		c.MatchReq()
		return
	}
	if bodyName == "pb.LoginResponse" {
		var res pb.LoginResponse
		proto.Unmarshal(down.Body, &res)
		if res.Code != constants.CODE_SUCCESS {
			gLog.Info("login err %+v msg %+v", res.Code, res.Msg)
			return
		}

		c.Uid = res.GetUid()
	}

	//if down.BodyName == "pb.LoginResponse" {
	//	var res pb.LoadProgressRes
	//	res.Unmarshal(down.Body)
	//
	//	isLoad := false
	//	for _, v := range res.GetData() {
	//		if v.Uid == c.Uid {
	//			isLoad = true
	//		}
	//	}
	//
	//	if isLoad {
	//		c.IsLoadProgress = true
	//		if c.IsHall {
	//			c.MatchReq()
	//		}
	//	}
	//
	//	return
	//}
	//
	//if down.BodyName == "pb.StartMatchRes" {
	//	messageMatch.AddSuccess(c.getReportSeq(down.Seq))
	//}
	//
	//if down.BodyName == "pb.MatchOverRes" {
	//	var res pb.MatchOverRes
	//	res.Unmarshal(down.Body)
	//	//ylog.Info(fmt.Sprintf("match res %v ", res))
	//	if res.GetCode() == 0 {
	//		messageMatch.AddSuccess(c.getReportSeq(down.Seq))
	//	} else {
	//		messageMatch.AddFail(c.getReportSeq(down.Seq))
	//	}
	//
	//	c.LeaveRoom()
	//
	//	c.IsEnterRoom = false
	//	c.IsLoadProgress = false
	//
	//	time.Sleep(time.Duration(1) * time.Second)
	//
	//	c.EndGameRoom(res)
	//
	//	c.IsEnterRoom = true
	//	return
	//}
	//
	//if down.BodyName == "pb.TestRes" {
	//	//ylog.Info(fmt.Sprintf("TestRes seq %v", down.Seq))
	//	var res pb.TestRes
	//	res.Unmarshal(down.Body)
	//	ylog.Info(fmt.Sprintf("uid %v test down seq %v", c.Uid, down.Seq))
	//	if res.GetCode() == 0 {
	//		messageTest.AddSuccess(c.getReportSeq(down.Seq))
	//	} else {
	//		messageTest.AddFail(c.getReportSeq(down.Seq))
	//	}
	//}
	//
	//if down.BodyName == "pb.DBTestRes" {
	//	//ylog.Info(fmt.Sprintf("DBTestRes seq %v", down.Seq))
	//	var res pb.DBTestRes
	//	res.Unmarshal(down.Body)
	//	if res.GetCode() == 0 {
	//		messageDBTest.AddSuccess(c.getReportSeq(down.Seq))
	//	} else {
	//		messageDBTest.AddFail(c.getReportSeq(down.Seq))
	//	}
	//}
	//
	//if down.BodyName == "pb.GameOverRes" {
	//	c.IsGameOver = true
	//	/* var res pb.GameOverRes
	//	res.Unmarshal(down.Body)
	//	ylog.Info(fmt.Sprintf("game over %v", res)) */
	//}
}

// Connect ...
func (c *MessageClient) Connect() error {
	conn := &websocket.Connection{
		Module:       "greeter",
		OnDownstream: c.Receive,
		Log:          gLog,
	}

	err := conn.Dial(c.URL)
	if err != nil {
		return fmt.Errorf("dial fail: %v", err)
	}

	c.Conn = conn
	return nil
}

func (c *MessageClient) HeartBeat() {
	req := &pb.HeartbeatRequest{
		Timestamp: time.Now().Unix(),
	}
	bytes, _ := proto.Marshal(req)
	head := &pb.ClientCommonHead{
		Pid:       c.Uid,
		Sn:        gCounter.GetIncrementValue(),
		ProtoName: proto.MessageName(req),
		Timestamp: time.Now().Unix(),
	}

	request := &pb.ClientCommonRequest{
		Head: head,
		Body: bytes,
	}

	err := c.Conn.Upstream(request)
	if err != nil {
		fmt.Println("message client send ping upstream fail: %v", err)
	}

}

func (c *MessageClient) SendEndHall() {
	if c.IsEnterRoom {
		return
	}

	gLog.Info(fmt.Sprintf("uid %v enter Hall", c.Uid))
	enterRoomReq := &pb.ClientLoginHallRequest{}
	bytes, _ := proto.Marshal(enterRoomReq)
	head := &pb.ClientCommonHead{
		Pid:       c.Uid,
		Sn:        gCounter.GetIncrementValue(),
		ProtoName: proto.MessageName(enterRoomReq),
		Timestamp: time.Now().Unix(),
	}

	request := &pb.ClientCommonRequest{
		Head: head,
		Body: bytes,
	}

	err := c.Conn.Upstream(request)
	if err != nil {
		fmt.Println("message client send ping upstream fail: %v", err)
	}

	messageEnterHall.AddRequest(c.getReportSeq(head.GetSn()))

	c.IsHall = true

	return
}

func (c *MessageClient) MatchReq() {
	if !c.IsHall {
		return
	}

	if c.IsMatch {
		return
	}

	c.IsMatch = true

	gLog.Info(fmt.Sprintf("uid %v match req", c.Uid))

	matchReq := &pb.ClientMatchRequest{}
	bytes, _ := proto.Marshal(matchReq)
	head := &pb.ClientCommonHead{
		Pid:       c.Uid,
		Sn:        gCounter.GetIncrementValue(),
		ProtoName: proto.MessageName(matchReq),
		Timestamp: time.Now().Unix(),
	}

	request := &pb.ClientCommonRequest{
		Head: head,
		Body: bytes,
	}

	err := c.Conn.Upstream(request)
	if err != nil {
		fmt.Println("message client send ping upstream fail: %v", err)
	}

	messageEnterHall.AddRequest(c.getReportSeq(head.GetSn()))
}

//func (c *MessageClient) EndGameRoom(matchRes pb.MatchOverRes) {
//	if c.IsEnterRoom {
//		return
//	}
//
//	ylog.Info(fmt.Sprintf("uid %v enter room", c.Uid))
//
//	c.RoomId = matchRes.RoomId
//	enterRoomReq := &gmegateway.GMEInternalEnterRoom{
//		Uid:      c.Uid,
//		UserData: []byte(matchRes.GetOpt()),
//	}
//
//	enterreq, _ := enterRoomReq.Marshal()
//	up := &gmegateway.Upstream{
//		Path:        g_protoPath,
//		ContentType: mimetype.ApplicationProtobuf,
//		BodyName:    proto.MessageName(enterRoomReq),
//		Body:        enterreq,
//		RoomId:      c.RoomId,
//	}
//
//	err := c.Conn.Upstream(up)
//	if err != nil {
//		fmt.Println("message client send ping upstream fail: %v", err)
//	}
//
//	messageEnterRoom.AddRequest(c.getReportSeq(up.Seq))
//}
//
//func (c *MessageClient) LeaveRoom() {
//	leaveRoomReq := &gmegateway.GMEInternalLeaveRoom{
//		Uid:      c.Uid,
//		UserData: nil,
//	}
//
//	leavereq, _ := leaveRoomReq.Marshal()
//
//	up := &gmegateway.Upstream{
//		Path:        g_protoPath,
//		ContentType: mimetype.ApplicationProtobuf,
//		BodyName:    proto.MessageName(leaveRoomReq),
//		Body:        leavereq,
//		RoomId:      c.RoomId,
//	}
//
//	err := c.Conn.Upstream(up)
//	if err != nil {
//		fmt.Println("message client send ping upstream fail: %v", err)
//	}
//}
//
//func (c *MessageClient) LoadProgress() {
//	if !c.IsEnterRoom {
//		return
//	}
//
//	ylog.Info(fmt.Sprintf("uid %v loadprogress", c.Uid))
//
//	if c.IsLoadProgress {
//		return
//	}
//
//	c.IsLoadProgress = true
//	loadProgressReq := &pb.LoadProgressReq{
//		Progress: 100,
//		Sn:       "111",
//	}
//
//	req, _ := loadProgressReq.Marshal()
//	up := &gmegateway.Upstream{
//		Path:        g_protoPath,
//		ContentType: mimetype.ApplicationProtobuf,
//		BodyName:    "pb.LoadProgressReq",
//		Body:        req,
//		RoomId:      c.RoomId,
//	}
//
//	err := c.Conn.Upstream(up)
//	if err != nil {
//		fmt.Println("message client send ping upstream fail: %v", err)
//		return
//	}
//
//	return
//}
//

//func (c *MessageClient) Test() {
//	if !c.IsEnterRoom {
//		return
//	}
//
//	if c.IsHall {
//		return
//	}

//	testReq := &pb.TestReq{Msg: RandString()}
//	test_req, _ := testReq.Marshal()
//	up := &gmegateway.Upstream{
//		Path:        g_protoPath,
//		ContentType: mimetype.ApplicationProtobuf,
//		BodyName:    proto.MessageName(testReq),
//		Body:        test_req,
//		RoomId:      c.RoomId,
//	}
//
//	err := c.Conn.Upstream(up)
//	if err != nil {
//		fmt.Println("message client send ping upstream fail: %v", err)
//	}
//
//	ylog.Info(fmt.Sprintf("uid %v test up seq %v", c.Uid, up.Seq))
//	messageTest.AddRequest(c.getReportSeq(up.Seq))
//}
//
//func (c *MessageClient) DbTest() {
//	if !c.IsEnterRoom {
//		return
//	}
//
//	if c.IsHall {
//		return
//	}
//
//	if c.ExcuteNum%10 == 0 {
//		dbTestReq := &pb.DBTestReq{}
//		dbtest_req, _ := dbTestReq.Marshal()
//		up := &gmegateway.Upstream{
//			Path:        g_protoPath,
//			ContentType: mimetype.ApplicationProtobuf,
//			BodyName:    proto.MessageName(dbTestReq),
//			Body:        dbtest_req,
//			RoomId:      c.RoomId,
//		}
//		err := c.Conn.Upstream(up)
//		if err != nil {
//			fmt.Println("message client send ping upstream fail: %v", err)
//		}
//
//		messageDBTest.AddRequest(c.getReportSeq(up.Seq))
//	}
//}

// Execute ...
func (c *MessageClient) Execute() error {
	c.ExcuteNum++
	if c.IsFirst {
		c.SendEndHall()
		c.IsFirst = false
	}

	c.HeartBeat()

	if !c.IsLoadProgress {
		return nil
	}

	if c.IsHall {
		return nil
	}

	if c.IsGameOver {
		return nil
	}

	//c.Test()
	//c.DbTest()

	return nil
}

// Close ...
func (c *MessageClient) Close() error {
	c.Conn.Close()
	return nil
}
