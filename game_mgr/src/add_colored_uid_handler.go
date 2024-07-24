package game_mgr

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"io"
	"net/http"
)

func (self *HttpService) AddColoredUidHandler(body []byte, w http.ResponseWriter) {
	var request domain.ColoredUidRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		self.Log.Info(" AddColoredUidHandler err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}
	self.Log.Info("AddColoredUidHandler request %+v", request)

	setKey := "COLORED_UID_SET_KEY" + request.GameId

	redisKey := "COLORED_UID_KEY" + request.GameId

	self.RedisDao.SAdd(setKey, request.UidList)

	self.RedisDao.Set(redisKey, request.Activity, 0)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
