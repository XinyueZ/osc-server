package osc

import (
	"bytes"
	"fmt"

	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strconv"
)


type Error string

func (e Error) Error() string {
	return string(e)
}

func makeHeader(r *http.Request, cookie string, length int) {
	r.Header.Add("Content-Type", API_RESTYPE)
	r.Header.Add("Content-Length", strconv.Itoa(length))
	r.Header.Add("Host", OSC)
	r.Header.Add("Connection", KEEP_ALIVE)
	r.Header.Add("Cookie", cookie)
}

func printHeader(r *http.Request) {
	header := r.Header
	for k, v := range header {
		fmt.Println("k:", k, "v:", v)
	}
}

type OsChina struct {
	XMLName xml.Name `xml:"oschina"`
	User    User     `xml:"user"`
}

type User struct {
	Uid  string `xml:"uid"`
	Name string `xml:"name"`
}

func init() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/tweetList", handleTweetList)
	http.HandleFunc("/tweetPub", handleTweetPub)
}



func main() {
	chUser := make(chan *User)
	chLogin := make(chan string)
	chTweetList := make(chan string)
	chTweetPub := make(chan string)
	defer func() {
		if e := recover(); e != nil {
			close(chUser)
			close(chLogin)
			close(chTweetList)
			close(chTweetPub)
		}
	}()
	go login(ACCOUNT, PWD, chLogin, chUser)
	cookie := <-chLogin //Got user session.
	cookie = "oscid=" + cookie
	puser := <-chUser
	if cookie != "" {
		fmt.Println(cookie)
		fmt.Println(puser.Uid)
		go printTweetList(puser, cookie, 1, chTweetList)
		tweetListContent := <-chTweetList
		if tweetListContent != "" {
			fmt.Println(tweetListContent)
			//Just a randem msg
			msgRandem := "做不了爱就不爱。。。"
			go pubTweet(puser, cookie, msgRandem, chTweetPub)
			pubContent := <-chTweetPub
			if pubContent != "" {
				fmt.Println(pubContent)

			}
		}
	}
}

func login(account string, password string, cookieCh chan string, userCh chan *User) {
	fmt.Println("Login.")
	client := new(http.Client)
	body := fmt.Sprintf(LOGIN_SCHEME, account, password)
	url := LOGIN_VALIDATE_HTTP
	fmt.Println(url)
	if r, e := http.NewRequest(POST, url, bytes.NewBufferString(body)); e == nil {
		makeHeader(r, "", len(body))
		if resp, e := client.Do(r); e == nil {
			fmt.Println(resp.Status)
			var cookie string = ""
			if resp != nil {
				defer resp.Body.Close()
			}
			if bytes, err := ioutil.ReadAll(resp.Body); err == nil {
				var posc OsChina
				if err := xml.Unmarshal(bytes, &posc); err == nil {
					for _, v := range resp.Cookies() {
						if v.Value != "" {
							cookie = v.Value
							break
						}
					}
					cookieCh <- cookie
					userCh <- &(posc.User)
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

func printTweetList(puser *User, cookie string, page int, ch chan string) {
	fmt.Println("Get Tweet-List.")
	client := new(http.Client)
	url := fmt.Sprintf(TWEET_LIST, puser.Uid, page)
	fmt.Println(url)
	if r, e := http.NewRequest(GET, url, nil); e == nil {
		makeHeader(r, cookie, 0)
		if resp, e := client.Do(r); e == nil {
			fmt.Println(resp.Status)
			if resp != nil {
				defer resp.Body.Close()
			}
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				ch <- string(bytes)
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

func pubTweet(puser *User, cookie string, msg string, ch chan string) {
	fmt.Printf("Pub Tweet: %s\n", msg)
	client := new(http.Client)
	url := TWEET_PUB
	fmt.Println(url)
	body := fmt.Sprintf(TWEET_PUB_SCHEME, puser.Uid, msg)
	fmt.Println(body)
	if r, e := http.NewRequest(POST, url, bytes.NewBufferString(body)); e == nil {
		makeHeader(r, cookie, len(body))
		printHeader(r)
		if resp, e := client.Do(r); e == nil {
			fmt.Println(resp.Status)
			if resp != nil {
				defer resp.Body.Close()
			}
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				ch <- string(bytes)
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
