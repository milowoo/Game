package constants

const (
	GM_CODE_SUBJECT           = "gm.code.subj."
	UCENTER_APPLY_UID_SUBJECT = "ucenter.apply.uid"
	LOGIN_HALL_SUBJECT        = "login.hall.req.subj"
	MATCH_SUBJECT             = "match.req.subj"
	MATCH_OVER_SUBJECT        = "match.over.subject"
)

func GetGmCodeSubject(gameId string) string {
	return GM_CODE_SUBJECT + gameId
}
