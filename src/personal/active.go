package personal

import (
	"common"

	"encoding/json"
)

type ActivesList struct {
	Notice       common.Notice `json:"notice"`
	ActivesArray []Active      `json:"activelist"`
}

func (self ActivesList) StringActivesArray() (s string) {
	json, _ := json.Marshal(&self.ActivesArray)
	s = string(json)
	return
}

func (self ActivesList) StringNotice() (s string) {
	json, _ := json.Marshal(&self.Notice)
	s = string(json)
	return
}

type Active struct {
	Id            int    `json:"id"`
	Portrait      string `json:"portrait"`
	Author        string `json:"author"`
	AuthorId      int    `json:"authorid"`
	Catalog       int    `json:"catalog"`
	AppClient     int    `json:"appClient"`
	ObjectId      int    `json:"objectId"`
	ObjectType    int    `json:"objectType"`
	ObjectCatalog int    `json:"objectCatalog"`
	ObjectTitle   string `json:"objectTitle"`
	Url           string `json:"url"`
	Message       string `json:"message"`
	TweetImage    string `json:"tweetImage"`
	CommentCount  int    `json:"commentCount"`
	PubDate       string `json:"pubDate"`
	ObjectReply   Reply  `json:"objectReply"`
}

type Reply struct {
	ObjectName string `json:"objectName"`
	ObjectBody string `json:"objectBody"`
}
