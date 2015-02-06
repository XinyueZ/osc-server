package osc

import (
	"encoding/xml"
)

type UserInfo struct {
	XMLName xml.Name `xml:"oschina"`
	User    User     `xml:"user"`
}

type User struct {
	Uid  string `xml:"uid"`
	Name string `xml:"name"`
}

type TweetList struct {
	XMLName     xml.Name    `xml:"oschina"`
	TweetsArray TweetsArray `xml:"tweets"`
}

type TweetsArray struct {
	Tweets []Tweet `xml:"tweet"`
}

type Tweet struct {
	Id           int    `xml:"id"`
	Portrait     string `xml:"portrait"`
	Author       string `xml:"author"`
	AuthorId     int    `xml:"authorid"`
	Body         string `xml:"body"`
	CommentCount int    `xml:"commentCount"`
	PubDate      string `xml:"pubDate"`
	ImgSmall     string `xml:"imgSmall"`
	ImgBig       string `xml:"imgBig"`
}

type ResultInfo struct {
	XMLName xml.Name `xml:"oschina"`
	Result  Result   `xml:"result"`
}

type Result struct {
	Code    int `xml:"errorCode"`
	Message string `xml:"errorMessage"`
}
