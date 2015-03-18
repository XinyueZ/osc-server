package personal

import (
	"common"
	"encoding/json"

	"appengine"
	"appengine/urlfetch"

	"bytes"
	"fmt"

	"io/ioutil"
	"net/http"
)

func AllFriendList(cxt appengine.Context, session string, access_token string,  relation int,  total int, ch chan *FriendsList)  (pFriendList  *FriendsList) {
	if total > 0 {
		page := 1
		go  FriendList(cxt, session, access_token, page, relation, ch)
		pFriendList = <- ch
		for total > len( pFriendList.FriendsArray ) {
			page++
			go  FriendList(cxt, session, access_token, page, relation, ch)
			pMoreFriendsList := <-ch
			pFriendList.FriendsArray  = append(pFriendList.FriendsArray, pMoreFriendsList.FriendsArray...)
		}
	} else {
		pFriendList = nil
	}
	return
}


func FriendList(cxt appengine.Context, session string, access_token string, page int, relation int, ch chan *FriendsList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.PERSONAL_FRIENDS_LIST_SCHEME,page, relation, access_token)
	//fmt.Fprintf(w, `%s\n`, body)
	if r, e := http.NewRequest(common.POST, common.PERSONAL_FRIENDS_LIST_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		//fmt.Fprintf(w, `oscid=%s\n`, session)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pFriendsList := new(FriendsList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				//fmt.Fprintf(w, `%s\n`, string(bytes))
				if err := json.Unmarshal(bytes, pFriendsList); err == nil {
					ch <- pFriendsList
				} else {
					panic(e)
				}
			} else {
				panic(e)
			}
		} else {
			panic(e)
		}
	} else {
		panic(e)
	}
}
