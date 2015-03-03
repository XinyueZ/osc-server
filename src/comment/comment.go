package comment

import (
	"common"
	"encoding/json"
)

type CommentList struct {
	Notice   common.Notice `json:"notice"`
	Comments []Comment     `json:"commentList"`
}

func (self CommentList) StringNotice() (s string) {
	json, _ := json.Marshal(&self.Notice)
	s = string(json)
	return
}

func (self CommentList) StringCommentArray() (s string) {
	json, _ := json.Marshal(&self.Comments)
	s = string(json)
	return
}

type Comment struct {
	Content         string             `json:"content"`
	Id              int                `json:"id"`
	PubDate         string             `json:"pubDate"`
	ClientType      int                `json:"client_type"`
	CommentAuthor   string             `json:"commentAuthor"`
	CommentAuthorId int                `json:"commentAuthorId"`
	CommentPortrait string             `json:"commentPortrait"`
	Refers          []CommentListRefer `json:"refers"`
	Replies         []CommentListReply `json:"replies"`
}

type CommentListRefer struct {
	Refertitle string `json:"refertitle"`
	Referbody  string `json:"referbody"`
}

type CommentListReply struct {
	Rauthor   string `json:"rauthor"`
	RpubDate  string `json:"rpubDate"`
	RauthorId int    `json:"rauthorId"`
	Rcontent  string `json:"rcontent"`
}
