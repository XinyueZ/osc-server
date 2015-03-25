package personal

import (
	"common"

	"appengine"
	"appengine/urlfetch"

	"bytes"
	"encoding/xml"
	"fmt"

	"io/ioutil"
	"net/http"
)

func MyFriendList(cxt appengine.Context, session string, uid int, relation int, page int, ch chan *FriendsList) (pFriendList *FriendsList) {
	total := 3
	if total > 0 {
		go FriendList(cxt, session, uid, page, relation, total, ch)
		pFriendList = <-ch
	} else {
		pFriendList = nil
	}
	return
}

func HisFriendList(cxt appengine.Context, session string, uid int, friend int, relation int, page int, ch chan *FriendsList) (pFriendList *FriendsList) {
	total := 3
	if total > 0 {
		go FriendList(cxt, session, friend, page, relation, total, ch)
		pFriendList = <-ch
	} else {
		pFriendList = nil
	}
	return
}



func AllMyFriendList(cxt appengine.Context, session string, uid int, relation int, ch chan *FriendsList) (pFriendList *FriendsList) {
	chMyInfo := make(chan *MyInfo)
	go MyInformation(cxt, session, uid, chMyInfo)
	pMyInfo := <-chMyInfo

	if pMyInfo == nil {
		return
	}

	total := 0
	switch relation {
	case 0:
		total = pMyInfo.User.Fans
	case 1:
		total = pMyInfo.User.Follow
	}

	if total > 0 {
		go FriendList(cxt, session, uid, 0, relation, total, ch)
		pFriendList = <-ch
	} else {
		pFriendList = nil
	}
	return
}


func AllHisFriendList(cxt appengine.Context, session string, uid int, friend int, relation int, ch chan *FriendsList) (pFriendList *FriendsList) {
	chInfo := make(chan *UserInfo)
	go UserInformation(cxt, session, uid, friend, chInfo)
	pInfo := <-chInfo

	if pInfo == nil {
		return
	}

	total := 0
	switch relation {
	case 0:
		total = pInfo.User.Fans
	case 1:
		total = pInfo.User.Follow
	}

	if total > 0 {
		go FriendList(cxt, session, friend, 0, relation, total, ch)
		pFriendList = <-ch
	} else {
		pFriendList = nil
	}
	return
}

func FriendList(cxt appengine.Context, session string, uid int, page int, relation int, total int, ch chan *FriendsList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.PERSONAL_FRIENDS_LIST_SCHEME, uid, page, relation, total)
	if r, e := http.NewRequest(common.POST, common.PERSONAL_FRIENDS_LIST_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pFriendsList := new(FriendsList)
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
