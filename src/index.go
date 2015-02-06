package osc

import (
	"appengine"
	"appengine/datastore"

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
	http.HandleFunc("/myTweetList", handleMyTweetList)
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

	s := fmt.Sprintf(`{"status":%d, "user":{"name":"%s", "uid":"%d", "expired":"%s"}}`, STATUS_OK, puser.Name, puser.Uid, expires.String())

	q := datastore.NewQuery("OnlineUser").Filter("Uid =", puser.Uid)
	clients := make([]OnlineUser, 0)
	keys, _ := q.GetAll(cxt, &clients)
	if len(clients) > 0 {
		//Delete old one if find a existed item.
		datastore.DeleteMulti(cxt, keys)
	}

	ponlineUser := &OnlineUser{puser.Uid, puser.Name, session}
	datastore.Put(cxt, datastore.NewIncompleteKey(cxt, "OnlineUser", nil), ponlineUser)
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
	i, _ := strconv.Atoi(uid)
	session := getSession(cxt, i) //Get user-session
	p, _ := strconv.Atoi(page)

	go printTweetList(cxt, 0, session, p, chTweetList)
	tweets := <-chTweetList
	showTweetList(w, r, tweets, 0, p)
}

func showTweetList(w http.ResponseWriter, r *http.Request, tweets []Tweet, uid int, page int) {
		tweetsJson := ""
		for _, tw := range tweets {
			tw.Body = strings.Replace(tw.Body, `"`, "'", -1)
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

func handleMyTweetList(w http.ResponseWriter, r *http.Request) {
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
	i, _ := strconv.Atoi(uid)
	session := getSession(cxt, i) //Get user-session
	p, _ := strconv.Atoi(page)

	go printTweetList(cxt, i, session, p, chTweetList)
	tweets := <-chTweetList
	showTweetList(w, r, tweets, i, p)
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
	i, _ := strconv.Atoi(uid)
	session := getSession(cxt, i) //Get user-session

	go pubTweet(cxt, i, session, msg, chTweetPub)
	pubRet := <-chTweetPub
	s := fmt.Sprintf(`{"status":%d, "result":{"code":%d, "msg":"%s"}}`, STATUS_OK, pubRet.Code, pubRet.Message)
	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}
