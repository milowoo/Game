package handler

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/internal"
	"game_mgr/src/pb"
	"github.com/golang/protobuf/proto"
	"io"
	"net/http"
	"strconv"
	"time"
)

func GmCodeHandler(body []byte, w http.ResponseWriter) {
	log := internal.GLog
	var request domain.GmCode
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Info(" gmCode err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	testMatch("88888")
	//testMatchResponse("88888")

	//testLoginHall("23232")

	//验签

	//builder, _ := proto.Marshal(&pb.GmCodeRequest{
	//	Code:     request.Code,
	//	GameId:   request.GameId,
	//	Uid:      request.Uid,
	//	Opt:      request.Opt,
	//	Operator: request.Operator,
	//})
	//
	//err = self.NatsPool.Publish(constants.GetGmCodeSubject(request.GameId), builder)
	//
	//if err != nil {
	//	self.Log.Error("GmCodeHandler err %+v", err)
	//	httpRes := domain.Response{Code: constants.PUBLIC_SUBJECT_ERROR, Msg: "request game err", Data: ""}
	//	buf, _ := json.Marshal(httpRes)
	//	io.WriteString(w, string(buf))
	//	return
	//}

	//test only

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
	return
}

func testMatch(uid string) {
	matchReq := &pb.MatchRequest{
		GameId:         "test4",
		Uid:            uid,
		Score:          "1",
		TimeStamp:      strconv.FormatInt(time.Now().UnixMilli(), 10),
		ReceiveSubject: constants.MATCH_OVER_SUBJECT,
		Opt:            "112",
	}

	request, _ := proto.Marshal(matchReq)
	var response interface{}
	err := internal.NatsPool.Request(constants.MATCH_SUBJECT, string(request), &response, 3*time.Second)
	if err != nil {
		internal.GLog.Error("testMatch gameId %+v uid %+v match err %+v", err)
		return
	}
	internal.GLog.Info("testMatch response %+v", response)
	dataMap := response.(map[string]interface{})
	resBytes := []byte(dataMap["data"].(string))
	var res pb.MatchResponse
	err = proto.Unmarshal(resBytes, &res)
	if err != nil {
		internal.GLog.Error("testMatch err %+v", err)
		return
	}

	internal.GLog.Info("testMatch %+v", res)

	return
}

func testMatchResponse(uid string) {
	gameId := "test4"
	roomId := gameId + "_room_" + uid
	hostIp := "127.0.0.1"

	//通知
	response := &pb.MatchOverRes{
		Code:      constants.CODE_SUCCESS,
		Msg:       "Success",
		GameId:    gameId,
		Uid:       uid,
		RoomId:    roomId,
		Timestamp: strconv.FormatInt(time.Now().UnixMilli(), 10),
		UidList:   make([]*pb.MatchData, 0),
		AiUidList: make([]*pb.MatchData, 0),
		GameIp:    hostIp,
	}

	res, _ := proto.Marshal(response)
	internal.GLog.Info("AddMatchRequest response subject %+v", constants.MATCH_OVER_SUBJECT)
	internal.NatsPool.Publish(constants.MATCH_OVER_SUBJECT, string(res))
}

func testLoginHall(uid string) {
	request := pb.CreateHallRequest{
		Uid:    uid,
		GameId: "test4",
	}

	bytes, _ := proto.Marshal(&request)
	var response interface{}

	//发起进入大厅请求
	internal.GLog.Info("testLoginHall begin %+v", uid)
	internal.NatsPool.Request(constants.LOGIN_HALL_SUBJECT, string(bytes), &response, 3*time.Second)
	internal.GLog.Info("testLoginHall response %+v", response)

	dataMap := response.(map[string]interface{})
	resBytes := []byte(dataMap["data"].(string))
	var res pb.CreateHallResponse
	err := proto.Unmarshal(resBytes, &res)
	if err != nil {
		internal.GLog.Error("testLoginHall err %+v", err)
		return
	}

	internal.GLog.Info("testLoginHall %+v", res)

}

func testApplyUid() {
	appRequest := &pb.ApplyUidRequest{
		Pid: "23232323",
	}
	bytes, _ := proto.Marshal(appRequest)
	var reply interface{}
	internal.NatsPool.Request(constants.UCENTER_APPLY_UID_SUBJECT, string(bytes), &reply, 3*time.Second)
	internal.GLog.Info("reply is %+v", reply)
	var res pb.ApplyUidResponse
	proto.Unmarshal(bytes, &res)

	internal.GLog.Info("ApplyUid response uid [%+v] pid [%+v] code [%+v]",
		res.GetUid(), res.GetPid(), res.GetCode())
}
