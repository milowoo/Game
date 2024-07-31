package handler

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/internal"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"io"
	"net/http"
	"time"
)

func OnOffGameHandler(config *domain.NacosConfig, body []byte, w http.ResponseWriter) {
	log := internal.GLog
	var request domain.GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Info(" onOffGame err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	if request.Status != 1 && request.Status != 2 {
		log.Info(" onOffGame err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_STATUS, Msg: "invalid request status", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	dbData := internal.Mongo.GetGame(request.GameId)
	if dbData == nil {
		log.Info(" onOffGame err no exist gameId %+v err %+v ", request.GameId, err)
		httpRes := domain.Response{Code: constants.INVALID_STATUS, Msg: "game no exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	dbData.Status = request.Status
	dbData.UTime = time.Now()
	if len(request.Operator) > 0 {
		dbData.Operator = request.Operator
	}

	redisBuf, _ := json.Marshal(dbData)
	internal.RedisDao.Set(request.GameId, redisBuf, 0)

	internal.Mongo.SavaGame(request.GameId, dbData)
	//通告 gameId 变化
	publicResult, err := internal.IClient.PublishConfig(vo.ConfigParam{
		DataId:  config.GameDataId,
		Group:   config.GameGroup,
		Content: request.GameId,
	})

	log.Info("on off game public result %v err %+v", publicResult, err)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
