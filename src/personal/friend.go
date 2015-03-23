package personal

import (
	"encoding/xml"
	"fmt"
)

type FriendsArray []Friend

func (self FriendsArray) String() (s string) {
	s = "["
	i := 0
	m := len(self)
	for _, f := range self {
		t := fmt.Sprintf(
			`{"expertise":"%s", "name":"%s",  "userid":%d,  "gender":%d, "portrait" : "%s"}`,
			f.Expertise, f.Name, f.UserId, f.Gender, f.Portrait)
		s += t
		if i != m-1 {
			s += ","
		}
		i++
	}
	s += "]"

	return
}

type FriendsList struct {
	XMLName xml.Name `xml:"oschina"`
	Friends Friends  `xml:"friends"`
}

type Friends struct {
	FriendsArray FriendsArray `xml:"friend"`
}

func (self FriendsList) StringFriendsArray() (s string) {
	return self.Friends.FriendsArray.String()
}

type Friend struct {
	Expertise string `xml:"expertise"`
	Name      string `xml:"name"`
	UserId    int    `xml:"userid"`
	Gender    int    `xml:"gender"` //1-man|2-lady
	Portrait  string `xml:"portrait"`
}

func (self Friend) String() (s string) {
	json, _ := xml.Marshal(&self)
	s = string(json)
	return
}
