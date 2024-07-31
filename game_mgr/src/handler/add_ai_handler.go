package handler

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/internal"
	"io"
	"net/http"
)

func AddAiHandler(body []byte, w http.ResponseWriter) {
	log := internal.GLog
	var request domain.AiInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Info(" AddAiHandler err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	aiSetKey := "AI_UID_PID_SET_KEY"
	aiHKey := "AI_UID_PID_H_KEY"

	internal.RedisDao.SAdd(aiSetKey, request.Uid)
	internal.RedisDao.HSet(aiHKey, request.Uid, request.Pid)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
