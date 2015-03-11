package comment

import (
	"common"

	"appengine"
	"appengine/urlfetch"

	"bytes"
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"
)

func TweetCommentReply(cxt appengine.Context, session string, access_token string, id int, content string, receiverId int, authorId int, replyId int, ch chan *common.Result) (pResult *common.Result) {
	go CommentReply(cxt, session, access_token, id, 3, content, receiverId, authorId, replyId, ch)
	pResult = <-ch
	return
}

func CommentReply(cxt appengine.Context, session string, access_token string, id int, catalog int, content string, receiverId int, authorId int, replyId int, ch chan *common.Result) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.COMMENT_REPLY_SCHEME, id, catalog, content, receiverId, authorId, replyId, access_token)
	if r, e := http.NewRequest(common.POST, common.COMMENT_REPLY_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, len(body))
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pRes := new(common.Result)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if err := json.Unmarshal(bytes, pRes); err == nil {
					ch <- pRes
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
