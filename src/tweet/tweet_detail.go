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

func TweetDetail(cxt appengine.Context, session string, access_token string, id int, ch chan *Tweet) {
	ch <- SyncTweetDetail(cxt, session, access_token, id)
}

func SyncTweetDetail(cxt appengine.Context, session string, access_token string, id int) (pTweet *Tweet) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.TWEET_DETAIL_SCHEME, id, access_token)
	if r, e := http.NewRequest(common.POST, common.TWEET_DETAIL_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, len(body))
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pTweet = new(Tweet)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := json.Unmarshal(bytes, pTweet); e == nil {
					return
				} else {
					pTweet = nil
					cxt.Errorf("Error but still going: %v", e)
				}
			} else {
				panic(e)
			}
		} else {
			pTweet = nil
			cxt.Errorf("Error but still going: %v", e)
		}
	} else {
		panic(e)
	}
	return
}
