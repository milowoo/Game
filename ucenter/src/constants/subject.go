package constants

const (
	LOGIN_HALL_SUBJECT         = "login.hall.req.subj"
	MATCH_SUBJECT              = "match.req.subj"
	CANCEL_MATCH_SUBJECT       = "cancel_match.req.subj"
	CREATE_ROOM_NOTICE_SUBJECT = "create.room.notice."
	GM_CODE_SUBJECT            = "gm.code.subj."
	UCENTER_APPLY_UID_SUBJECT  = "ucenter.apply.uid"
)

func GetCreateRoomNoticeSubject(gameId, groupId string) string {
	return CREATE_ROOM_NOTICE_SUBJECT + gameId + "." + groupId
}

func GetGmCodeSubject(gameId string) string {
	return GM_CODE_SUBJECT + gameId
}

func GetGameSubject(gameId string, ip string) string {
	return "game." + gameId + "." + ip
}
