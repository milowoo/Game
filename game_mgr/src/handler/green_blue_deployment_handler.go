package handler

import (
	"encoding/json"
	"game_mgr/src"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"io"
	"net/http"
	"time"
)

func GreenBlueDeploymentHandler(self *game_mgr.HttpService, body []byte, w http.ResponseWriter) {
	var request domain.GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		self.Log.Info(" greenBlueDeployment err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	if len(request.GroupName) < 1 {
		self.Log.Info(" greenBlueDeployment err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_GROUP_NAME, Msg: "invalid request GroupName", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	dbData, err := self.Mongo.GetGame(request.GameId)
	if err == nil {
		self.Log.Info(" greenBlueDeployment err no exist gameId %+v err %+v ", request.GameId, err)
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
	self.RedisDao.Set(request.GameId, redisBuf, 0)

	self.Mongo.SavaGame(request.GameId, &dbData)

	//通告 gameId 变化
	self.IClient.PublishConfig(vo.ConfigParam{
		DataId:  self.GameConfig.NacosConfig.GameDataId,
		Group:   self.GameConfig.NacosConfig.GameGroup,
		Content: request.GameId,
	})

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
