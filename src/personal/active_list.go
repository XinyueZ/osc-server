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

func LastTweetActiveList(cxt appengine.Context, session string, access_token string, user int, page int, ch chan *ActivesList) (pActivesList *ActivesList) {
	go TweetActiveList(cxt, session, access_token, user, page, ch)
	pActivesList = <-ch
	atMoment := pActivesList.Notice.ReferCount
	if atMoment > 0 { //Only last new referes will be shown on client.
		pActivesList.ActivesArray = pActivesList.ActivesArray[:(atMoment)]
	} else {
		pActivesList = nil
	}
	return
}

func LastCommentActiveList(cxt appengine.Context, session string, access_token string, user int, page int, ch chan *ActivesList) (pActivesList *ActivesList) {
	go CommentsActiveList(cxt, session, access_token, user, page, ch)
	pActivesList = <-ch
	atMoment := pActivesList.Notice.ReplyCount
	if atMoment > 0 { //Only last new replies will be shown on client.
		pActivesList.ActivesArray = pActivesList.ActivesArray[:(atMoment)]
	} else {
		pActivesList = nil
	}
	return
}

func TweetActiveList(cxt appengine.Context, session string, access_token string, user int, page int, ch chan *ActivesList) {
	Actives(cxt, session, access_token, user, 2, page, ch)
}

func CommentsActiveList(cxt appengine.Context, session string, access_token string, user int, page int, ch chan *ActivesList) {
	Actives(cxt, session, access_token, user, 3, page, ch)
}

func Actives(cxt appengine.Context, session string, access_token string, user int, catalog int, page int, ch chan *ActivesList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.ACTIVE_LIST_SCHEME, catalog, user, page, access_token)
	//fmt.Fprintf(w, `%s\n`, body)
	if r, e := http.NewRequest(common.POST, common.ACTIVE_LIST_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		//fmt.Fprintf(w, `oscid=%s\n`, session)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pActivesList := new(ActivesList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				//fmt.Fprintf(w, `%s\n`, string(bytes))
				if err := json.Unmarshal(bytes, pActivesList); err == nil {
					ch <- pActivesList
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
