package comment

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/urlfetch"

	"common"
)

func TweetCommentList(cxt appengine.Context, session string, access_token string, id int, page int, ch chan *CommentList) {
	Comments(cxt, session, access_token, id, 3, page, ch)
}

func Comments(cxt appengine.Context, session string, access_token string, id int, catalog int, page int, ch chan *CommentList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.COMMENT_LIST_SCHEME, catalog, id, page, access_token)
	//fmt.Fprintf(w, `%s\n`, body)
	if r, e := http.NewRequest(common.POST, common.COMMENT_LIST_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		//fmt.Fprintf(w, `oscid=%s\n`, session)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pCommentList := new(CommentList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				//fmt.Fprintf(w, `%s\n`, string(bytes))
				if err := json.Unmarshal(bytes, pCommentList); err == nil {
					ch <- pCommentList
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
