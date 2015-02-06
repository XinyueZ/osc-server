package osc

import (
	"appengine"

	"fmt"
	"strings"

	"net/http"
	"strconv"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func init() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/tweetList", handleTweetList)
	http.HandleFunc("/tweetPub", handleTweetPub)
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chUser := make(chan *User)
	chLogin := make(chan *http.Cookie)
	defer func() {
		if e := recover(); e != nil {
			close(chUser)
			close(chLogin)
			fmt.Sprintf(`{"status":%d}`, STATUS_ERR)
		}
	}()
	args := r.URL.Query()
	account := args[ACCOUNT][0]
	pwd := args[PWD][0]
	go login(cxt, account, pwd, chLogin, chUser)

	//Get cookie.
	cookie := <-chLogin
	session := cookie.Value   //Got user session.
	expires := cookie.Expires //Time-up for session.
	//Get user-info
	puser := <-chUser

	s := fmt.Sprintf(`{"status":%d, "user":{"name":"%s", "uid":"%s", "expired":"%s"}}`, STATUS_OK, puser.Name, puser.Uid, expires.String())
	cookies := [...]*http.Cookie{
		&http.Cookie{
			Name:   "oscid",
			Value:  session,
			Path:   "/",
		},
	}
	for _, cookie := range cookies {
		http.SetCookie(w, cookie) //cookie to client
	}
	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}

func handleTweetList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetList := make(chan []Tweet)
	defer func() {
		if e := recover(); e != nil {
			close(chTweetList)
			fmt.Sprintf(`{"status":%d}`, STATUS_ERR)
		}
	}()
	args := r.URL.Query()
	uid := args[UID][0]   //Get user-id
	page := args[PAGE][0] //Which page
	cookies := r.Cookies() //Session in cookies passt
	session := cookies[0].Value //Get user-session
	i, _ := strconv.Atoi(uid)
	p, _ := strconv.Atoi(page)
	go printTweetList(cxt, i, session, p, chTweetList)
	tweets := <-chTweetList
	//<-chTweetList
	tweetsJson := ""
	for _, tw := range tweets {
		body := fmt.Sprintf(`{"id":%d, "pubDate":"%s", "body":"%s", "author":"%s", "authorid":%d, "imgSmall":"%s" , "commentCount":%d, "imgBig":"%s", "portrait":"%s"},`,
			tw.Id, tw.PubDate, tw.Body, tw.Author, tw.AuthorId, tw.ImgSmall, tw.CommentCount, tw.ImgBig, tw.Portrait)
		tweetsJson += body
	}
	tweetsJson = strings.Replace(tweetsJson, "<![CDATA[", "", -1)
	tweetsJson = strings.Replace(tweetsJson, "]]>", "", -1)
	tweetsJson = tweetsJson[:len(tweetsJson)-1] //Rmv last ","
	s := fmt.Sprintf(`{"status":%d, "tweets":[%s]}`, STATUS_OK, tweetsJson)
	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}

func handleTweetPub(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetPub := make(chan Result)
	defer func() {
		if e := recover(); e != nil {
			close(chTweetPub)
			fmt.Sprintf(`{"status":%d}`, STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	uid := args[UID][0] //Get user-id
	msg := args[MSG][0] //What to tweet
	cookies := r.Cookies() //Session in cookies passt
	session := cookies[0].Value //Get user-session
	i, _ := strconv.Atoi(uid)
	go pubTweet(cxt, i, session, msg, chTweetPub)
	pubRet := <-chTweetPub
	s := fmt.Sprintf(`{"status":%d, "result":{"code":%d, "msg":"%s"}}`, STATUS_OK, pubRet.Code, pubRet.Message)
	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}
