package game_mgr

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"io"
	"net/http"
	"time"
)

/*
增加游戏
*/
func (self *HttpService) AddGameHandler(body []byte, w http.ResponseWriter) {
	var request domain.GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		self.Log.Info(" addGame err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	if len(request.GameId) < 1 {
		self.Log.Info(" addGame err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_GAME_ID, Msg: "invalid request group data", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}
	if len(request.GroupName) < 1 {
		self.Log.Info(" addGame err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_GROUP_NAME, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	request.CTime = time.Now()
	request.UTime = time.Now()

	//判断是否存在这个game 信息
	result := self.Mongo.GetGame(request.GameId)
	if result != nil {
		self.Log.Info(" addGame err have exist gameId %+v err %+v ", request.GameId, err)
		httpRes := domain.Response{Code: constants.GAME_IS_EXIST, Msg: "game have exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	self.Log.Info("result %+v", result)

	//判断是否有记录
	data, _ := self.RedisDao.EXISTS(request.GameId)
	if data {
		self.Log.Info(" addGame err have exist gameId %+v err %+v ", request.GameId, err)
		httpRes := domain.Response{Code: constants.GAME_IS_EXIST, Msg: "game have exist", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	redisBuf, _ := json.Marshal(request)
	self.RedisDao.Set(request.GameId, redisBuf, 0)

	self.RedisDao.SAdd("ALL_GAME_ID", request.GameId)

	self.Mongo.InsertGame(&request)

	//通告 gameId 变化
	publicResult, err := self.IClient.PublishConfig(vo.ConfigParam{
		DataId:  self.GameConfig.NacosConfig.GameDataId,
		Group:   self.GameConfig.NacosConfig.GameGroup,
		Content: request.GameId,
	})

	self.Log.Info("PublishConfig result %+v err %+v", publicResult, err)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
