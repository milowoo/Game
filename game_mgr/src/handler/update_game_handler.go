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

func UpdateGameHanler(self *game_mgr.HttpService, body []byte, w http.ResponseWriter) {
	var request domain.GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		self.Log.Info(" addGame err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	//获取 game 信息
	dbData, err := self.Mongo.GetGame(request.GameId)
	if err == nil {
		self.Log.Info(" addGame err no exist gameId %+v err %+v ", request.GameId, err)
		httpRes := domain.Response{Code: constants.GAME_NOT_EXIST, Msg: "game no exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	if request.GameTime != 0 {
		dbData.GameTime = request.GameTime
	}
	if request.Type != 0 {
		dbData.Type = request.Type
	}

	if len(request.Name) > 1 {
		dbData.Name = request.Name
	}

	if request.MatchTime != 0 {
		dbData.MatchTime = request.MatchTime
	}

	dbData.UTime = time.Now()
	if len(request.Operator) > 0 {
		dbData.Operator = request.Operator
	}

	//记录 redis
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
