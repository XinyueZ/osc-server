package common

import "encoding/json"

type Notice struct {
	ReplyCount int `json:"replyCount"`
	MsgCount   int `json:"msgCount"`
	FansCount  int `json:"fansCount"`
	ReferCount int `json:"referCount"`
}

type Result struct {
	Code     string    `json:"error"`
	Relation int    `json:"relation"` 
}

func (self *Result) String() (s string) {
	json, _ := json.Marshal(&self)
	s = string(json)
	return
}
