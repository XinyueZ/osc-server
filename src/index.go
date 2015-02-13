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

	"common"
	"tweet"
	"user"
)

type Error string

func (e Error) Error() string {
	return string(e)
}

func init() {
	http.HandleFunc("/login", handleLogin)
	http.HandleFunc("/tweetList", handleTweetList)
	http.HandleFunc("/myTweetList", handleMyTweetList)
	http.HandleFunc("/hotspotTweetList", handleHotspotTweetList)
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
	chLogin := make(chan *user.Logined)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleLogin: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	base64Text, _ := ioutil.ReadAll(r.Body)
	if len(base64Text) == 0 {
		panic("Needs login account info.")
		return
	}

	ciphertext := decodeBase64(string(base64Text))
	decrypted := make([]byte, len(ciphertext))
	decryptAESCFB(decrypted, ciphertext, common.KEY, common.IV)
	plainText := string(decrypted)

	data := strings.Split(plainText, "&")
	account := strings.TrimSpace((strings.Split(data[0], "="))[1])
	pwd := strings.TrimSpace((strings.Split(data[1], "="))[1])

	pUser := user.NewOscUser(account, pwd, common.APP_ID, common.APP_SEC)
	go pUser.Login(cxt, chLogin)

	//Get cookie.
	pLogined := <-chLogin
	session := pLogined.Cookie.Value    //Got user session.
	expires := pLogined.Token.ExpiresIn //Time-up for session.
	uid := pLogined.Token.UID
	access_token := pLogined.Token.AccessToken
	s := fmt.Sprintf(`{"status":%d, "user":{"uid":%d, "expired":%d}}`, common.STATUS_OK, uid, expires)

	w.Header().Set("Set-Cookie", "oscid="+session+";access_tokien="+access_token)
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//--------------------------------------------------------------------------------
//List of tweets, my tweets, all tweets(uid==0) etc.
//--------------------------------------------------------------------------------

//Show user's tweets.
func handleMyTweetList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetList := make(chan *tweet.TweetList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleMyTweetList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	uid := args[common.UID][0]              //Get user-id
	page := args[common.PAGE][0]            //Which page
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	i, _ := strconv.Atoi(uid)
	p, _ := strconv.Atoi(page)

	go tweet.PrintTweetList(cxt, i, session, access_token, p, chTweetList)
	tweet.ShowTweetList(w, r, <-chTweetList, i, p)
}

//Show all tweets.
func handleTweetList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetList := make(chan *tweet.TweetList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	page := args[common.PAGE][0]            //Which page
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	p, _ := strconv.Atoi(page)
	go tweet.PrintTweetList(cxt, 0, session, access_token, p, chTweetList)
	tweet.ShowTweetList(w, r, <-chTweetList, 0, p)
}

//Show hotspot tweets.
func handleHotspotTweetList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetList := make(chan *tweet.TweetList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	page := args[common.PAGE][0]     //Which page
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	p, _ := strconv.Atoi(page)
	go tweet.PrintTweetList(cxt, -1, session, access_token, p, chTweetList)
	tweet.ShowTweetList(w, r, <-chTweetList, 0, p)
}

//--------------------------------------------------------------------------------

//Publish a tweets.
func handleTweetPub(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetPub := make(chan *common.Result)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetPub: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	uid := args[common.UID][0]              //Get user-id
	msg := args[common.MSG][0]              //What to tweet
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	i, _ := strconv.Atoi(uid)

	go tweet.PubTweet(cxt, i, session, access_token, msg, chTweetPub)
	pRes := <-chTweetPub
	code, _ := strconv.Atoi(pRes.Code)
	message := pRes.Message
	s := fmt.Sprintf(`{"status":%d, "result":{"code":%d, "msg":"%s"}}`, common.STATUS_OK, code, message)
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}
