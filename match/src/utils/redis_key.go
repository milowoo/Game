package utils

func GetMatchUserKey(gameId, uid string) string {
	return gameId + ".match.user." + uid
}
