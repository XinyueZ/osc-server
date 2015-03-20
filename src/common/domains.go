package common

import (
	"appengine"
	"appengine/urlfetch"

	"bytes"
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

type Notice struct {
	ReplyCount int `json:"replyCount"`
	MsgCount   int `json:"msgCount"`
	FansCount  int `json:"fansCount"`
	ReferCount int `json:"referCount"`
}

type Result struct {
	Code     string `json:"error"`
	Relation int    `json:"relation"`
}

func (self *Result) String() (s string) {
	json, _ := json.Marshal(&self)
	s = string(json)
	return
}

func ClearAtNotice(cxt appengine.Context, session string, access_token string, ch chan *Result) (pResult *Result) {
	go ClearNotice(cxt, session, access_token, 1, ch)
	pResult = <-ch
	return
}

func ClearCommentsNotice(cxt appengine.Context, session string, access_token string, ch chan *Result) (pResult *Result) {
	go ClearNotice(cxt, session, access_token, 3, ch)
	pResult = <-ch
	return
}

//Clear notice with different type.
func ClearNotice(cxt appengine.Context, session string, access_token string, typ int, ch chan *Result) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(CLEAR_NOTICE_SCHEME, typ, access_token)
	if r, e := http.NewRequest(POST, CLEAR_NOTICE_URL, bytes.NewBufferString(body)); e == nil {
		MakeHeader(r, "oscid="+session, len(body))
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pResult := new(Result)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := json.Unmarshal(bytes, pResult); e == nil {
					ch <- pResult
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
