package common

type Notice struct {
	ReplyCount int `json:"replyCount"`
	MsgCount   int `json:"msgCount"`
	FansCount  int `json:"fansCount"`
	ReferCount int `json:"referCount"`
}

type Result struct {
	Code    string `json:"error"`
	Message string `json:"error_description"`
}
