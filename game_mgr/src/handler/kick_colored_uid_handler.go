package handler

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/internal"
	"io"
	"net/http"
)

func KickColoredUidHandler(body []byte, w http.ResponseWriter) {
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
	log.Info("kickColoredUidHandler request %+v", request)

	setKey := "COLORED_UID_SET_KEY" + request.GameId

	internal.RedisDao.SREM(setKey, request.UidList)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
