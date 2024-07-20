package domain

type ColoredUidRequest struct {
	GameId   string   `json:"gameId,omitempty"`
	Activity string   `json:"activity,omitempty"`
	UidList  []string `json:"UidList,omitempty"`
}
