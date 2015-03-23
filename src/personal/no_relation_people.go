package personal

import (
	"appengine"
)

//Get no relation people, they are friends other my friends.
func GetNoRelationPeople(cxt appengine.Context, session string, access_token string, uid int) (s string) {
	chFansList := make(chan *FriendsList)
	chFollowList := make(chan *FriendsList)
	chUserInfo := make(chan *UserInfo)
	chOtherFriendsFans := make(chan *FriendsList)
	chOtherFriendsFollow := make(chan *FriendsList)

	//Get all friends, inc. fans, followers.
	pFans := AllMyFriendList(cxt, session, uid, 0, chFansList)
	pfollow := AllMyFriendList(cxt, session, uid, 1, chFollowList)

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
	for _, friend := range allFriends {
		go AllHisFriendList(cxt, session, uid, friend.UserId, 0, chOtherFriendsFans)
		go AllHisFriendList(cxt, session, uid, friend.UserId, 1, chOtherFriendsFollow)
	}

	//Combine all friends of other user.
	friendsOfOther := make([]Friend, 0)
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
