package tweet

import (
	"common"
	"encoding/json"
)

type TweetsList struct {
	Notice      common.Notice `json:"notice"`
	TweetsArray []Tweet       `json:"tweetlist"`
}

func (self TweetsList) StringTweetsArray() (s string) {
	json, _ := json.Marshal(&self.TweetsArray)
	s = string(json)
	return
}

func (self TweetsList) StringNotice() (s string) {
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
