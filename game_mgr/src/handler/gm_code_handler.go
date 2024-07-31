package handler

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/internal"
	"io"
	"net/http"
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
	internal.NatsPool.Publish(constants.GetGmCodeSubject(request.GameId), "test")

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
	return
}
