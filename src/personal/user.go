package personal

import (
	"common"
	"fmt"

	"encoding/json"
)

type UserInfo struct {
	Notice        common.Notice `json:"notice"`
	Uid           int           `json:"uid"`
	Name          string        `json:"name"`
	Ident         string        `json:"ident"`
	Province      string        `json:"province"`
	City          string        `json:"city"`
	Platforms     []string      `json:"platforms"`
	Expertise     []string      `json:"expertise"`
	Portrait      string        `json:"portrait"`
	Gender        int           `json:"gender"`   //1-man, 2,famle
	Relation      int           `json:"relation"` //1-has been focused, 2-focused eachother, 3-no
	JoinTime      string        `json:"joinTime"`
	LastLoginTime string        `json:"lastLoginTime"`
}

func (self UserInfo) String() (s string) {
	json, _ := json.Marshal(&self)
	s = string(json)
	s = fmt.Sprintf(
		`{"uid":%d, "name":"%s", "ident" : "%s","province":"%s", "city" : "%s","platforms":"%s", "expertise" : "%s","portrait":"%s", "gender" : %d,"relation":%d}`,
		self.Uid,
		self.Name,
		self.Ident,
		self.Province, self.City,
		convert(self.Platforms),
		convert(self.Expertise),
		self.Portrait,
		self.Gender,
		self.Relation)
	return
}

func (self UserInfo) StringNotice() (s string) {
	json, _ := json.Marshal(&self.Notice)
	s = string(json)
	return
}

func convert(a []string) (s string) {
	s = ""
	if a != nil && len(a) > 0 {
		for _, v := range a {
			s += (v + " ")
		}
	}
	return
}

type MyInfo struct {
	Notice         common.Notice `json:"notice"`
	Uid            int           `json:"uid"`
	Name           string        `json:"name"`
	Ident          string        `json:"ident"`
	Province       string        `json:"province"`
	City           string        `json:"city"`
	Platforms      []string      `json:"platforms"`
	Expertise      []string      `json:"expertise"`
	Portrait       string        `json:"portrait"`
	Gender         int           `json:"gender"`   //1-man, 2,famle
	Relation       int           `json:"relation"` //1-has been focused, 2-focused eachother, 3-no
	JoinTime       string        `json:"joinTime"`
	LastLoginTime  string        `json:"lastLoginTime"`
	FansCount      int           `json:"fansCount"`
	FavoriteCount  int           `json:"favoriteCount"`
	FollowersCount int           `json:"followersCount"`
}

func (self MyInfo) String() (s string) {
	json, _ := json.Marshal(&self)
	s = string(json)
	s = fmt.Sprintf(
		`{"uid":%d, "name":"%s", "ident" : "%s","province":"%s", "city" : "%s","platforms":"%s", "expertise" : "%s","portrait":"%s", "gender" : %d,"relation":%d, "fansCount":%d, "favoriteCount":%d,"followersCount":%d}`,
		self.Uid,
		self.Name,
		self.Ident,
		self.Province, self.City,
		convert(self.Platforms),
		convert(self.Expertise),
		self.Portrait,
		self.Gender,
		self.Relation,
		self.FansCount,
		self.FavoriteCount,
		self.FollowersCount)
	return
}

func (self MyInfo) StringNotice() (s string) {
	json, _ := json.Marshal(&self.Notice)
	s = string(json)
	return
}
