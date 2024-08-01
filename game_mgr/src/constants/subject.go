package constants

const (
	GM_CODE_SUBJECT           = "gm.code.subj."
	UCENTER_APPLY_UID_SUBJECT = "ucenter.apply.uid"
)

func GetGmCodeSubject(gameId string) string {
	return GM_CODE_SUBJECT + gameId
}
