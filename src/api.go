package osc

import (
	"appengine"
	"appengine/urlfetch"

	"bytes"
	"fmt"

	"encoding/xml"
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

func login(cxt appengine.Context, account string, password string, cookieCh chan *http.Cookie, userCh chan *User) {
	fmt.Println("Login.")
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(LOGIN_SCHEME, account, password)
	url := LOGIN_VALIDATE_HTTP
	fmt.Println(url)
	if r, e := http.NewRequest(POST, url, bytes.NewBufferString(body)); e == nil {
		makeHeader(r, "", len(body))
		if resp, e := client.Do(r); e == nil {
			fmt.Println(resp.Status)
			var cookie *http.Cookie
			if resp != nil {
				defer resp.Body.Close()
			}
			if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
				var osc UserInfo
				if err := xml.Unmarshal(bytes, &osc); err == nil {
					for _, v := range resp.Cookies() {
						if v.Value != "" {
							cookie = v
							break
						}
					}
					cookieCh <- cookie
					userCh <- &(osc.User)
				} else {
					panic(err)
				}

			} else {
				panic(err)
			}
		} else {
			panic(e)
		}
	} else {
		panic(e)
	}
}

func printTweetList(cxt appengine.Context, uid int, session string, page int, ch chan []Tweet) {
	fmt.Println("Get Tweet-List.")
	client := urlfetch.Client(cxt)
	url := fmt.Sprintf(TWEET_LIST, uid, page)
	fmt.Println(url)
	body := fmt.Sprintf(TWEET_LIST_SCHEME, uid, page)
	if r, e := http.NewRequest(POST, url, bytes.NewBufferString(body)); e == nil {
		makeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			fmt.Println(resp.Status)
			if resp != nil {
				defer resp.Body.Close()
			}
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				var tweetList TweetList
				if err := xml.Unmarshal(bytes, &tweetList); err == nil {
					tweets := tweetList.TweetsArray.Tweets
					ch <- tweets
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

func pubTweet(cxt appengine.Context, uid int, session string, msg string, ch chan Result) {
	fmt.Println("Get Tweet-Pub.")
	client := urlfetch.Client(cxt)
	url := TWEET_PUB
	fmt.Println(url)
	body := fmt.Sprintf(TWEET_PUB_SCHEME, uid, msg)
	if r, e := http.NewRequest(POST, url, bytes.NewBufferString(body)); e == nil {
		makeHeader(r, "oscid="+session, len(body))
		if resp, e := client.Do(r); e == nil {
			fmt.Println(resp.Status)
			if resp != nil {
				defer resp.Body.Close()
			}
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				var ri ResultInfo
				if err := xml.Unmarshal(bytes, &ri); err == nil {
					ch <- ri.Result
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
