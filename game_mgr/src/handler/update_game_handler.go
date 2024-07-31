package handler

import (
	"encoding/json"
	"game_mgr/src/constants"
	"game_mgr/src/domain"
	"game_mgr/src/internal"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"io"
	"net/http"
)

func UpdateGameHandler(config *domain.NacosConfig, body []byte, w http.ResponseWriter) {
	log := internal.GLog
	var request domain.GameInfo
	err := json.Unmarshal(body, &request)
	if err != nil {
		log.Info(" addGame err %v ", err)
		httpRes := domain.Response{Code: constants.INVALID_BODY, Msg: "invalid request body", Data: ""}
		buf, _ := json.Marshal(httpRes)
		io.WriteString(w, string(buf))
		return
	}

	//获取 game 信息
	dbData := internal.Mongo.GetGame(request.GameId)
	if dbData == nil {
		log.Info(" addGame err no exist gameId %+v err %+v ", request.GameId, err)
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

	if len(request.Operator) > 0 {
		dbData.Operator = request.Operator
	}

	//记录 redis
	redisBuf, _ := json.Marshal(dbData)
	internal.RedisDao.Set(request.GameId, redisBuf, 0)

	internal.Mongo.SavaGame(request.GameId, dbData)

	//通告 gameId 变化
	publicResult, err := internal.IClient.PublishConfig(vo.ConfigParam{
		DataId:  config.GameDataId,
		Group:   config.GameGroup,
		Content: request.GameId,
	})

	log.Info("update public DataId %+v group %+v result %v err %+v",
		config.GameDataId,
		config.GameGroup,
		publicResult, err)

	httpRes := domain.Response{Code: constants.CODE_SUCCESS, Msg: "success", Data: ""}
	buf, _ := json.Marshal(httpRes)
	io.WriteString(w, string(buf))
}
