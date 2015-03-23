package personal

import (
	"common"
	"encoding/xml"
	"fmt"
)

type UserInfo struct {
	XMLName xml.Name     `xml:"oschina"`
	User    UserInfoData `xml:"user"`
}

type UserInfoData struct {
	Uid           int    `xml:"uid"`
	Name          string `xml:"name"`
	From          string `xml:"from"`
	Platforms     string `xml:"devplatform"`
	Expertise     string `xml:"expertise"`
	Portrait      string `xml:"portrait"`
	Gender        string `xml:"gender"`   //1-man, 2,famle
	Relation      int    `xml:"relation"` //1-has been focused, 2-focused eachother, 3-no
	Score         int    `xml:"score"`
	Fans          int    `xml:"fans"`
	Follow        int    `xml:"followers"`
	JoinTime      string `xml:"jointime"`
	LastLoginTime string `xml:"latestonline"`
}

func (self UserInfo) String() (s string) {
	s = fmt.Sprintf(
		`{"uid":%d, "name":"%s",  "from":"%s",  "platforms":"%s", "expertise" : "%s","portrait":"%s", "gender" : %d,"relation":%d, "fans":%d, "follow":%d}`,
		self.User.Uid,
		self.User.Name,
		self.User.From,
		self.User.Platforms,
		self.User.Expertise,
		self.User.Portrait,
		genderConver(self.User.Gender),
		self.User.Relation,
		self.User.Fans,
		self.User.Follow)
	return
}

type MyInfo struct {
	XMLName xml.Name          `xml:"oschina"`
	User    UserInfoData      `xml:"user"`
	Notice  common.NoticeData `xml:"notice"`
}

func (self MyInfo) String() (s string) {
	s = fmt.Sprintf(
		`{"uid":%d, "name":"%s",  "from":"%s",  "platforms":"%s", "expertise" : "%s","portrait":"%s", "gender" : %d,"relation":%d, "fans":%d, "follow":%d}`,
		self.User.Uid,
		self.User.Name,
		self.User.From,
		self.User.Platforms,
		self.User.Expertise,
		self.User.Portrait,
		genderConver(self.User.Gender),
		self.User.Relation,
		self.User.Fans,
		self.User.Follow)
	return
}

func genderConver(gender string) (n int) {
	if gender == "男" {
		n = 1
	} else {
		n = 2
	}
	return
}
