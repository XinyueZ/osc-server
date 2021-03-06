package tweet

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

func TweetPub(cxt appengine.Context, session string, access_token string, msg string, ch chan *common.Result) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.TWEET_PUB_SCHEME, access_token, msg)
	if r, e := http.NewRequest(common.POST, common.TWEET_PUB_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, len(body))
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pRes := new(common.Result)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := json.Unmarshal(bytes, pRes); e == nil {
					ch <- pRes
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
