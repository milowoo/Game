package game_mgr

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/pb"
	"github.com/golang/protobuf/proto"
	"io"
	"net/http"
	"reflect"
	"time"
)

func (self *HttpService) GmCodeHandler(body []byte, w http.ResponseWriter) {
	var response interface{}

	request := &pb.ApplyUidRequest{
		Pid: "0000fasdf0oadsd",
	}

	req, _ := proto.Marshal(request)
	err := self.NatsPool.Request(constants.UCENTER_APPLY_UID_SUBJECT, string(req), &response, 3*time.Second)
	if err != nil {
		self.Log.Info("GmCodeHandler err %+v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid reply", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	dataType := reflect.TypeOf(response)
	self.Log.Info("GmCodeHandler response type %+v response %v", dataType, response)
	dataMap := response.(map[string]interface{})
	bytes := []byte(dataMap["data"].(string))
	var res pb.ApplyUidResponse
	proto.Unmarshal(bytes, &res)

	self.Log.Info("GmCodeHandler response uid [%+v] pid [%+v] code [%+v]",
		res.GetUid(), res.GetPid(), res.GetCode())

	//var appResponse pb.ApplyUidRequest
	//proto.Unmarshal([]byte(res), &appResponse)
	//self.Log.Info("GmCodeHandler appResponse %+v", appResponse)

	//var request domain.GmCode
	//err := json.Unmarshal(body, &request)
	//if err != nil {
	//	self.Log.Info(" gmCode err %v ", err)
	//	httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
	//	buf, _ := json.Marshal(httpRes)
	//	io.WriteString(w, string(buf))
	//	return
	//}
	//
	////验签
	//
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
	//
	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
	return
}
