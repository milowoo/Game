package handler

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/internal"
	"io"
	"net/http"
)

func DeleteColoredUidHandler(body []byte, w http.ResponseWriter) {
	log := internal.GLog
	var request domain.ColoredUidRequest
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Info(" AddColoredUidHandler err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}
	log.Info("AddColoredUidHandler request %+v", request)

	setKey := "COLORED_UID_SET_KEY" + request.GameId

	redisKey := "COLORED_UID_KEY" + request.GameId

	internal.RedisDao.Del(setKey)

	internal.RedisDao.Del(redisKey)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
