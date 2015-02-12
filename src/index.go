package osc

import (
	"appengine"

	"fmt"
	"strings"

	"io/ioutil"
	"net/http"
	"strconv"

	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
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

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

func decryptAESCFB(dst, src, key, iv []byte) {
	aesBlockDecrypter, _ := aes.NewCipher(key)
	aesDecrypter := cipher.NewCFBDecrypter(aesBlockDecrypter, iv)
	aesDecrypter.XORKeyStream(dst, src)
}

//Login a user and store oschina session-id to cookie.
//Use the cookie to access rest APIs.
func handleLogin(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chLogin := make(chan *Logined)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleLogin: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, STATUS_ERR)
		}
	}()

	base64Text, _ := ioutil.ReadAll(r.Body)
	if len(base64Text) == 0 {
		panic("Needs login account info.")
		return
	}

	ciphertext := decodeBase64(string(base64Text))
	decrypted := make([]byte, len(ciphertext))
	decryptAESCFB(decrypted, ciphertext, KEY, IV)
	plainText := string(decrypted)

	data := strings.Split(plainText, "&")
	account := strings.TrimSpace((strings.Split(data[0], "="))[1])
	pwd := strings.TrimSpace((strings.Split(data[1], "="))[1])

	pUser := newOscUser(account, pwd, APP_ID, APP_SEC)
	go pUser.login(cxt, chLogin)

	//Get cookie.
	pLogined := <-chLogin
	session := pLogined.Cookie.Value    //Got user session.
	expires := pLogined.Token.ExpiresIn //Time-up for session.
	uid := pLogined.Token.UID
	access_token := pLogined.Token.AccessToken
	s := fmt.Sprintf(`{"status":%d, "user":{"uid":%d, "expired":%d}}`, STATUS_OK, uid, expires)

	w.Header().Set("Set-Cookie", "oscid="+session+";access_tokien="+access_token)
	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}

//--------------------------------------------------------------------------------
//List of tweets, my tweets, all tweets(uid==0) etc.
//--------------------------------------------------------------------------------

//Show tweets in array(slice) by paging whit a uid.
//When uid==0 then show all. Show also tweets of user self by his uid.
func showTweetList(w http.ResponseWriter, r *http.Request, pTweetList *TweetList, uid int, page int) {
	s := fmt.Sprintf(`{"status":%d, "tweets":%s}`, STATUS_OK, pTweetList.StringTweetsArray())
	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Show user's tweets.
func handleMyTweetList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetList := make(chan *TweetList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleMyTweetList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	uid := args[UID][0]              //Get user-id
	page := args[PAGE][0]            //Which page
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	i, _ := strconv.Atoi(uid)
	p, _ := strconv.Atoi(page)

	go printTweetList(cxt,  i, session, access_token, p, chTweetList)
	showTweetList(w, r, <-chTweetList, i, p)
}

//Show all tweets.
func handleTweetList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetList := make(chan *TweetList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	page := args[PAGE][0]            //Which page
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	p, _ := strconv.Atoi(page)
	go printTweetList(cxt,  0, session, access_token, p, chTweetList)
	showTweetList(w, r, <-chTweetList, 0, p)
}

//--------------------------------------------------------------------------------

//Publish a tweets.
func handleTweetPub(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetPub := make(chan *Result)
	defer func() {
		if err := recover(); err != nil { 
			cxt.Errorf("handleTweetPub: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	uid := args[UID][0]              //Get user-id
	msg := args[MSG][0]              //What to tweet
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	i, _ := strconv.Atoi(uid)

	go pubTweet(cxt, i, session, access_token, msg, chTweetPub)
	pRes := <-chTweetPub
	code, _ := strconv.Atoi(pRes.Code)
	message := pRes.Message
	s := fmt.Sprintf(`{"status":%d, "result":{"code":%d, "msg":"%s"}}`, STATUS_OK, code, message)
	w.Header().Set("Content-Type", API_RESTYPE)
	fmt.Fprintf(w, s)
}
