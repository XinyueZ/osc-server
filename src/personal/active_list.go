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

func LastTweetActiveList(cxt appengine.Context, session string, access_token string, user int, page int, showMe bool, ch chan *ActivesList) (pActivesList *ActivesList) {
	go TweetActiveList(cxt, session, access_token, user, page, showMe, ch)
	pActivesList = <-ch
	atMoment := pActivesList.Notice.ReferCount
	if atMoment > 0 { //Only last new referes will be shown on client.
		pActivesList.ActivesArray = pActivesList.ActivesArray[:(atMoment)]
	} else {
		pActivesList = nil
	}
	return
}

func LastCommentActiveList(cxt appengine.Context, session string, access_token string, user int, page int, showMe bool, ch chan *ActivesList) (pActivesList *ActivesList) {
	go CommentsActiveList(cxt, session, access_token, user, page, showMe, ch)
	pActivesList = <-ch
	atMoment := pActivesList.Notice.ReplyCount
	if atMoment > 0 { //Only last new replies will be shown on client.
		pActivesList.ActivesArray = pActivesList.ActivesArray[:(atMoment)]
	} else {
		pActivesList = nil
	}
	return
}

func TweetActiveList(cxt appengine.Context, session string, access_token string, user int, page int, showMe bool, ch chan *ActivesList) {
	Actives(cxt, session, access_token, user, 2, page, showMe, ch)
}

func CommentsActiveList(cxt appengine.Context, session string, access_token string, user int, page int, showMe bool, ch chan *ActivesList) {
	Actives(cxt, session, access_token, user, 3, page, showMe, ch)
}

func Actives(cxt appengine.Context, session string, access_token string, user int, catalog int, page int, showMe bool, ch chan *ActivesList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.ACTIVE_LIST_SCHEME, catalog, user, page, access_token)
	if r, e := http.NewRequest(common.POST, common.ACTIVE_LIST_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pActivesList := new(ActivesList)    //Only for temp .
			pActivesListRet := new(ActivesList) //Real to return.
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if err := json.Unmarshal(bytes, pActivesList); err == nil {
					if !showMe {
						pActivesListRet.Notice = pActivesList.Notice
						pActivesListRet.ActivesArray = []Active{}
						for _, v := range pActivesList.ActivesArray {
							if v.AuthorId != user {
								pActivesListRet.ActivesArray = append(pActivesListRet.ActivesArray, v)
							}
						}
					} else {
						pActivesListRet = pActivesList
					}
					ch <- pActivesListRet
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
