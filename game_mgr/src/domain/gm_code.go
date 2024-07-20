package domain

type GmCode struct {
	Code     int32  `json:"code,omitempty"`
	GameId   string `json:"gameId,omitempty"`
	Uid      string `json:"uid,omitempty"`
	Opt      string `json:"opt,omitempty"`
	SignStr  string `json:"signStr,omitempty"`
	Operator string `json:"operator,omitempty"`
}
