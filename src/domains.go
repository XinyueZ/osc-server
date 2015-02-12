package osc

import (
	"encoding/json"
)

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

type TweetList struct {
	Notice      Notice  `json:"notice"`
	TweetsArray []Tweet `json:"tweetlist"`
}

func (self TweetList) StringTweetsArray() (s string) {
	json, _ := json.Marshal(&self.TweetsArray)
	s = string(json)
	return
}

func (self TweetList) StringNotice() (s string) {
	json, _ := json.Marshal(&self.Notice)
	s = string(json)
	return
}

type Tweet struct {
	Id           int    `json:"id"`
	Portrait     string `json:"portrait"`
	Author       string `json:"author"`
	AuthorId     int    `json:"authorid"`
	Body         string `json:"body"`
	CommentCount int    `json:"commentCount"`
	PubDate      string `json:"pubDate"`
	ImgSmall     string `json:"imgSmall"`
	ImgBig       string `json:"imgBig"`
}

func (self Tweet) String() (s string) {
	json, _ := json.Marshal(&self)
	s = string(json)
	return
}
