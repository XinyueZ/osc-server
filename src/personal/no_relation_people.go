package personal

import (
	"appengine"
	"appengine/urlfetch"

	"bytes"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"

	"common"
)

//In order to get friends-list of other people, we do not use openAPI directly.
//Instead openAPI we use APIs which was used in oschina mobile client.
const (
	POST                    = "POST"
	API_FRIENDS_URL         = "http://www.oschina.net/action/api/friends_list"
	API_FRIENDS_SCHEME      = "uid=%d&pageIndex=%d&relation=%d&pageSize=20"
	API_FRIENDS_INFO_URL    = "http://www.oschina.net/action/api/user_information"
	API_FRIENDS_INFO_SCHEME = "uid=%d&hisuid=%d&hisname=%s&pageIndex=0&pageSize=0"
)

type OtherFriendsList struct {
	XMLName xml.Name     `xml:"oschina"`
	Friends OtherFriends `xml:"friends"`
}

type OtherFriends struct {
	FriendsArray []OtherFriend `xml:"friend"`
}

type OtherFriend struct {
	Expertise string `xml:"expertise"`
	Name      string `xml:"name"`
	UserId    int    `xml:"userid"`
	Gender    int    `xml:"gender"` //1-man|2-lady
	Portrait  string `xml:"portrait"`
}

//Specical API to get friends of other people, it is not usage of openAPI.
func otherFriendList(cxt appengine.Context, session string, access_token string, uid int, page int, relation int, ch chan *OtherFriendsList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(API_FRIENDS_SCHEME, uid, page, relation)
	if r, e := http.NewRequest(POST, API_FRIENDS_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pFriendsList := new(OtherFriendsList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := xml.Unmarshal(bytes, pFriendsList); e == nil {
					ch <- pFriendsList
				} else {
					ch <- nil
					cxt.Errorf("Error but still going: %v", e)
				}
			} else {
				ch <- nil
				panic(e)
			}
		} else {
			ch <- nil
			cxt.Errorf("Error but still going: %v", e)
		}
	} else {
		ch <- nil
		panic(e)
	}
}



//Get no relation people, they are friends other my friends.
func GetNoRelationPeople(cxt appengine.Context, session string, access_token string, uid int) (s string) {
	chMyInfo := make(chan *MyInfo)
	chFansList := make(chan *FriendsList)
	chFollowList := make(chan *FriendsList)
	chUserInfo := make(chan *UserInfo)
	chOtherFriendsFans := make(chan *OtherFriendsList)
	chOtherFriendsFollow := make(chan *OtherFriendsList)

	//Get number of current friends.
	go MyInformation(cxt, session, uid, chMyInfo)
	pMyInfo := <-chMyInfo

	if pMyInfo == nil {
		return "null"
	}

	//Get all friends, inc. fans, followers.
	pFans := AllFriendList(cxt, session, access_token, 0, pMyInfo.User.Fans, chFansList)
	pfollow := AllFriendList(cxt, session, access_token, 1, pMyInfo.User.Follow, chFollowList)

	if pFans == nil && pfollow == nil {
		return "null"
	}

	//Combine all friends.
	allFriends := make([]Friend, 0)
	if pFans != nil && pfollow != nil {
		allFriends = append(pFans.FriendsArray, pfollow.FriendsArray...)
	} else if pFans != nil {
		allFriends = pFans.FriendsArray
	} else {
		allFriends = pfollow.FriendsArray
	}

	//Get all possible friends.
	for _, friend := range allFriends {
		go otherFriendList(cxt, session, access_token, friend.UserId, 0, 0, chOtherFriendsFans)
		go otherFriendList(cxt, session, access_token, friend.UserId, 0, 1, chOtherFriendsFollow)
	}

	//Combine all friends of other user.
	friendsOfOther := make([]OtherFriend, 0)
	for i := 0; i < len(allFriends); i++ {
		f := <-chOtherFriendsFans
		if f != nil {
			friendsOfOther = append(friendsOfOther, f.Friends.FriendsArray...)
		}
		f = <-chOtherFriendsFollow
		if f != nil {
			friendsOfOther = append(friendsOfOther, f.Friends.FriendsArray...)
		}
	}

	//Filter out people that have been added in my friends.
	availables := make([]OtherFriend, 0)
	for _, fo := range friendsOfOther {
		if fo.UserId != uid && !inList(&fo, pFans) && !inList(&fo, pfollow) && !inAvailables(&fo, availables[:]) {
			availables = append(availables, fo)
		}
	}

	l := len(availables)
	if l > 0 {
		//Make feeds for client to provide all user-information.
		for _, a := range availables {
			go UserInformation(cxt, session, uid, a.UserId, chUserInfo)
		}

		s = "["
		for i := 0; i < l; i++ {
			userInfo := <-chUserInfo
			if userInfo != nil {
				s += userInfo.String()
				s += ","

			}
		}
		s = s[:len(s)-1]
		s += "]"
	} else {
		s = "null"
	}
	return
}

//To know whether the author of an active in my current friends-list or not.
func inList(pFriendOther *OtherFriend, pFriends *FriendsList) bool {
	for _, f := range pFriends.FriendsArray {
		if f.UserId == pFriendOther.UserId {
			return true
		}
	}
	return false
}

func inAvailables(pFriendOther *OtherFriend, availables []OtherFriend) bool {
	for _, f := range availables {
		if f.UserId == pFriendOther.UserId {
			return true
		}
	}
	return false
}
