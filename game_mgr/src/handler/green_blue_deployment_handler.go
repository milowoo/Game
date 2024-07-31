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

func GreenBlueDeploymentHandler(config *domain.NacosConfig, body []byte, w http.ResponseWriter) {
	log := internal.GLog
	var request domain.GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Info(" greenBlueDeployment err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	if len(request.GroupName) < 1 {
		log.Info(" greenBlueDeployment err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_GROUP_NAME, Msg: "invalid request GroupName", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	dbData := internal.Mongo.GetGame(request.GameId)
	if dbData == nil {
		log.Info(" greenBlueDeployment err no exist gameId %+v err %+v ", request.GameId, err)
		httpRes := domain.Response{Code: constants.GAME_NOT_EXIST, Msg: "game no exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	dbData.GroupName = request.GroupName
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

	log.Info("green blue game public result %v err %+v", publicResult, err)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
