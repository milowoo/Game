package constants

const (
	//系统错误码
	CODE_SUCCESS         = 100
	INVALID_BODY         = 101
	PUBLIC_SUBJECT_ERROR = 102
	INVALID_GAME_ID      = 103
	GAME_IS_OFF          = 104
	SYSTEM_ERROR         = 105
	INVALID_REQUEST      = 106

	//  1000 - 2000  gateway
	INVALID_TIME       = 1001
	EXPIRED_TIME       = 1002
	PLAYER_IS_MATCHING = 1003
	PLAYER_NO_MATCHING = 1004
	PLAYER_IN_ROOM     = 1005

	//2000 - 3000    game-mgr
	INVALID_URL        = 2001 /* invalid url */
	INVALID_GROUP_NAME = 2002
	GAME_IS_EXIST      = 2004
	GAME_NOT_EXIST     = 2005
	INVALID_STATUS     = 2006

	//3000 - 4000 match

	// 4000 - 5000 ucenter

	// 5000 - 10000 game
)