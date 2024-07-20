package handler

import (
	"encoding/json"
	"game_mgr/src"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"io"
	"net/http"
)

func DeleteColoredUidHandler(self *game_mgr.HttpService, body []byte, w http.ResponseWriter) {
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

	self.RedisDao.Del(setKey)

	self.RedisDao.Del(redisKey)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
