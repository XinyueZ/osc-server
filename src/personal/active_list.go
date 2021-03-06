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

func LastTweetActiveList(cxt appengine.Context, session string, access_token string, user int, showMe bool, ch chan *ActivesList) (pActivesList *ActivesList) {
	page := 1
	go Actives(cxt, session, access_token, user, 2, page, ch)
	pActivesList = <-ch
	atMoment := pActivesList.Notice.ReferCount

	if atMoment > 0 {
		for atMoment > len(pActivesList.ActivesArray) {
			page++
			go Actives(cxt, session, access_token, user, 2, page, ch)
			pMoreActivesList := <-ch
			pActivesList.ActivesArray = append(pActivesList.ActivesArray, pMoreActivesList.ActivesArray...)
		}
		pActivesList.ActivesArray = pActivesList.ActivesArray[:(atMoment)]
		if !showMe {
			pActivesListRet := new(ActivesList)
			pActivesListRet.Notice = common.Notice{0, 0, 0, 0}
			pActivesListRet.ActivesArray = []Active{}
			for _, v := range pActivesList.ActivesArray {
				if v.AuthorId != user {
					pActivesListRet.ActivesArray = append(pActivesListRet.ActivesArray, v)
					pActivesListRet.Notice.ReferCount++
				}
			}
			pActivesList = pActivesListRet
		}
	} else {
		pActivesList = nil
	}
	return
}

func LastCommentActiveList(cxt appengine.Context, session string, access_token string, user int, showMe bool, ch chan *ActivesList) (pActivesList *ActivesList) {
	page := 1
	go Actives(cxt, session, access_token, user, 3, page, ch)
	pActivesList = <-ch
	atMoment := pActivesList.Notice.ReplyCount

	if atMoment > 0 {
		for atMoment > len(pActivesList.ActivesArray) {
			page++
			go Actives(cxt, session, access_token, user, 3, page, ch)
			pMoreActivesList := <-ch
			pActivesList.ActivesArray = append(pActivesList.ActivesArray, pMoreActivesList.ActivesArray...)
		}
		pActivesList.ActivesArray = pActivesList.ActivesArray[:(atMoment)]
		if !showMe {
			pActivesListRet := new(ActivesList)
			pActivesListRet.Notice = common.Notice{0, 0, 0, 0}
			pActivesListRet.ActivesArray = []Active{}
			for _, v := range pActivesList.ActivesArray {
				if v.AuthorId != user {
					pActivesListRet.ActivesArray = append(pActivesListRet.ActivesArray, v)
					pActivesListRet.Notice.ReplyCount++
				}
			}
			pActivesList = pActivesListRet
		}
	} else {
		pActivesList = nil
	}
	return
}

func TweetActiveList(cxt appengine.Context, session string, access_token string, user int, page int, showMe bool, ch chan *ActivesList) {
	Actives(cxt, session, access_token, user, 2, page, ch)
}

func CommentsActiveList(cxt appengine.Context, session string, access_token string, user int, page int, showMe bool, ch chan *ActivesList) {
	Actives(cxt, session, access_token, user, 3, page, ch)
}

func Actives(cxt appengine.Context, session string, access_token string, user int, catalog int, page int, ch chan *ActivesList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.ACTIVE_LIST_SCHEME, catalog, user, page, access_token)
	if r, e := http.NewRequest(common.POST, common.ACTIVE_LIST_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pActivesList := new(ActivesList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := json.Unmarshal(bytes, pActivesList); e == nil {
					ch <- pActivesList
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
