package constants

const (
	GM_CODE_SUBJECT = "gm.code.subj."
)

func GetGmCodeSubject(gameId string) string {
	return GM_CODE_SUBJECT + gameId
}
