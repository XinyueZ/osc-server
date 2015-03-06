package personal

import (
	"common" 

	"encoding/json"
)

type FriendsList struct {
	Notice       common.Notice `json:"notice"`
	FriendsArray []Friend      `json:"userList"`
}

func (self FriendsList) StringFriendsArray() (s string) {
	json, _ := json.Marshal(&self.FriendsArray)
	s = string(json)
	return
}

func (self FriendsList) StringNotice() (s string) {
	json, _ := json.Marshal(&self.Notice)
	s = string(json)
	return
}

type Friend struct {
	Expertise string `json:"expertise"`
	Name      string `json:"name"`
	UserId    int    `json:"userid"`
	Gender    int    `json:"gender"` //1-man|2-lady
	Portrait  string `json:"portrait"`
}

func (self Friend) String() (s string) {
	json, _ := json.Marshal(&self)
	s = string(json)
	return
}

