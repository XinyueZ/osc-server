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

type TweetList struct {
	Notice      common.Notice `json:"notice"`
	TweetsArray []Tweet       `json:"tweetlist"`
}

func (self TweetList) StringTweetsArray() (s string) {
	json, _ := json.Marshal(&self.TweetsArray)
	s = string(json)
	return
}

func (self TweetList) StringNotice() (s string) {
	json, _ := json.Marshal(&self.Notice)
	s = string(json)
	return
}

type Tweet struct {
	Id           int    `json:"id"`
	Portrait     string `json:"portrait"`
	Author       string `json:"author"`
	AuthorId     int    `json:"authorid"`
	Body         string `json:"body"`
	CommentCount int    `json:"commentCount"`
	PubDate      string `json:"pubDate"`
	ImgSmall     string `json:"imgSmall"`
	ImgBig       string `json:"imgBig"`
}

func (self Tweet) String() (s string) {
	json, _ := json.Marshal(&self)
	s = string(json)
	return
}

//Show tweets in array(slice) by paging whit a uid.
//When uid==0 then show all. Show also tweets of user self by his uid.
func ShowTweetList(w http.ResponseWriter, r *http.Request, pTweetList *TweetList, uid int, page int) {
	s := fmt.Sprintf(`{"status":%d, "tweets":%s}`, common.STATUS_OK, pTweetList.StringTweetsArray())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

func PrintTweetList(cxt appengine.Context, uid int, session string, access_token string, page int, ch chan *TweetList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.TWEET_LIST_SCHEME, uid, access_token, page)
	if r, e := http.NewRequest(common.POST, common.TWEET_LIST_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pTweetList := new(TweetList)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if err := json.Unmarshal(bytes, pTweetList); err == nil {
					ch <- pTweetList
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
