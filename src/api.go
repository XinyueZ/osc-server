package osc

import (
	"appengine"
	"appengine/urlfetch"

	"bytes"
	"fmt"

	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
)

func makeHeader(r *http.Request, cookie string, length int) {
	r.Header.Add("Content-Type", API_REQTYPE)
	r.Header.Add("Content-Length", strconv.Itoa(length))
	r.Header.Add("Host", OSC)
	r.Header.Add("Connection", KEEP_ALIVE)
	r.Header.Add("Cookie", cookie)
}

func printHeader(r *http.Request, w http.ResponseWriter) {
	header := r.Header
	for k, v := range header {
		fmt.Fprintf(w, "k:%s v:%s", k, v)
	}
}

func printTweetList(cxt appengine.Context,  uid int, session string, access_token string, page int, ch chan *TweetList) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(TWEET_LIST_SCHEME, uid, access_token, page)
	if r, e := http.NewRequest(POST, TWEET_LIST, bytes.NewBufferString(body)); e == nil {
		makeHeader(r, "oscid="+session, 0)
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

func pubTweet(cxt appengine.Context, uid int, session string, access_token string, msg string, ch chan *Result) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(TWEET_PUB_SCHEME, uid, access_token, msg)
	if r, e := http.NewRequest(POST, TWEET_PUB, bytes.NewBufferString(body)); e == nil {
		makeHeader(r, "oscid="+session, len(body))
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pRes := new(Result)
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
