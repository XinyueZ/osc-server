package personal

import (
	"appengine"
)

//Get no relation people, they are friends other my friends.
func GetNoRelationPeople(cxt appengine.Context, session string, access_token string,uid int, myPage int, friendPage int ) (s string) {
	chFansList := make(chan *FriendsList)
	chFollowList := make(chan *FriendsList)
	chUserInfo := make(chan *UserInfo)
	chOtherFriendsFans := make(chan *FriendsList)
	chOtherFriendsFollow := make(chan *FriendsList)

	//Get all friends, inc. fans, followers.
	pFans := MyFriendList(cxt, session, uid, 0, myPage, chFansList)
	pfollow := MyFriendList(cxt, session, uid, 1, myPage, chFollowList)

	if pFans == nil && pfollow == nil {
		return "null"
	}


	//Combine all friends.
	allFriends := make([]Friend, 0)
	if pFans != nil && pfollow != nil {
		allFriends = append(pFans.Friends.FriendsArray, pfollow.Friends.FriendsArray...)
	} else if pFans != nil {
		allFriends = pFans.Friends.FriendsArray
	} else {
		allFriends = pfollow.Friends.FriendsArray
	}

	//Get all possible friends.
	friendsOfOther := make([]Friend, 0)
	for _, friend := range allFriends {
		//friend := allFriends[i]
		pFans   := HisFriendList(cxt, session, uid, friend.UserId, 0, friendPage, chOtherFriendsFans)
		if pFans != nil {
			friendsOfOther = append(friendsOfOther, pFans.Friends.FriendsArray...)
		}
		pfollow := HisFriendList(cxt, session, uid, friend.UserId, 1, friendPage, chOtherFriendsFollow)
		if pfollow != nil {
			friendsOfOther = append(friendsOfOther, pfollow.Friends.FriendsArray...)
		}
	}


	//Filter out people that have been added in my friends.
	availables := make([]Friend, 0)
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
func inList(pFriendOther *Friend, pFriends *FriendsList) bool {
	for _, f := range pFriends.Friends.FriendsArray {
		if f.UserId == pFriendOther.UserId {
			return true
		}
	}
	return false
}

func inAvailables(pFriendOther *Friend, availables []Friend) bool {
	for _, f := range availables {
		if f.UserId == pFriendOther.UserId {
			return true
		}
	}
	return false
}
