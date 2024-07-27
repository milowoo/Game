package websocket

import (
	"net/url"
	"strconv"
	"time"
	"yalitest/util"
)

// URL ...
type URL struct {
	// 连接地址
	Addr string
	// 游戏名称 (game id)
	GameName string
	// 厂商AppKey
	AppKey string
	// uid
	UID string
}

// String ...
func (u *URL) String() string {
	tm := time.Now().Unix()

	rawURL, _ := url.Parse(u.Addr)

	v := url.Values{}
	v.Add("timestamp", strconv.FormatInt(tm, 10))
	v.Add("token", util.UUID())
	v.Add("pid", util.UUID())
	v.Add("gameId", u.GameName)
	rawURL.RawQuery = v.Encode()

	return rawURL.String()
}

// postDataInfo 连接信息
// Note: 数据结构来自匹配服务, 见 [游戏登录协议](http://kxd_3rd.kxd-pages.yy.com/kxd_3rd_doc/game/denglu.html)
type postDataInfo struct {
	// 游戏渠道
	ChannelID string `json:"channelid"`
	// 游戏名
	GameName string `json:"gameid"`
	// 用户信息
	Player postDataPlayerInfo `json:"player"`
}

// postDataPlayerInfo 用户信息
// Note: 数据结构来自匹配服务, 见 [游戏登录协议](http://kxd_3rd.kxd-pages.yy.com/kxd_3rd_doc/game/denglu.html)
type postDataPlayerInfo struct {
	// 用户ID
	UID string `json:"uid"`
	// 用户昵称
	Name string `json:"name"`
	// 用户头像URL
	AvatarURL string `json:"avatarurl"`
	// 用户性别
	Sex int64 `json:"sex"`
}
