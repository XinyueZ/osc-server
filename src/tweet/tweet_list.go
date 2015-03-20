package tweet

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

//Show tweets in array(slice) by paging whit a uid.
//When uid==0 then show all. Show also tweets of user self by his uid.
func ShowTweetList(w http.ResponseWriter, r *http.Request, pTweetsList *TweetsList, uid int, page int) {
	s := fmt.Sprintf(`{"status":%d, "tweets":%s}`, common.STATUS_OK, pTweetsList.StringTweetsArray())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

func TweetList(cxt appengine.Context, uid int, session string, access_token string, page int, ch chan *TweetsList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.TWEET_LIST_SCHEME, uid, access_token, page)
	if r, e := http.NewRequest(common.POST, common.TWEET_LIST_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pTweetsList := new(TweetsList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := json.Unmarshal(bytes, pTweetsList); e == nil {
					ch <- pTweetsList
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
