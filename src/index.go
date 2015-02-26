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
	"personal"
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
	http.HandleFunc("/friendsList", handleFriendsList)
	http.HandleFunc("/userInformation", handlePersonal)
	http.HandleFunc("/updateRelation", handleUpdateRelation)
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
	chTweetList := make(chan *tweet.TweetsList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleMyTweetList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	uid := args[common.UID][0]       //Get user-id
	page := args[common.PAGE][0]     //Which page
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	i, _ := strconv.Atoi(uid)
	p, _ := strconv.Atoi(page)

	go tweet.TweetList(cxt, i, session, access_token, p, chTweetList)
	tweet.ShowTweetList(w, r, <-chTweetList, i, p)
}

//Show all tweets.
func handleTweetList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetList := make(chan *tweet.TweetsList)
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
	go tweet.TweetList(cxt, 0, session, access_token, p, chTweetList)
	tweet.ShowTweetList(w, r, <-chTweetList, 0, p)
}

//Show hotspot tweets.
func handleHotspotTweetList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetList := make(chan *tweet.TweetsList)
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
	go tweet.TweetList(cxt, -1, session, access_token, p, chTweetList)
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
	uid := args[common.UID][0]       //Get user-id
	msg := args[common.MSG][0]       //What to tweet
	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	i, _ := strconv.Atoi(uid)

	go tweet.TweetPub(cxt, i, session, access_token, msg, chTweetPub)
	pRes := <-chTweetPub
	//code, _ := strconv.Atoi(pRes.Code) 
	s := fmt.Sprintf(`{"status":%d, "result":%s}`, common.STATUS_OK, pRes.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//--------------------------------------------------------------------------------
//List of friends, fans and who are focused,
//--------------------------------------------------------------------------------

//Get all friends
func handleFriendsList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chFansList := make(chan *personal.FriendsList)
	chFocusList := make(chan *personal.FriendsList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleFriendsList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	//0-fans|1-who are focused.
	go personal.FriendList(cxt, session, access_token, 0, chFansList)
	go personal.FriendList(cxt, session, access_token, 1, chFocusList)
	pFans := <-chFansList
	pFocus := <-chFocusList
	s := fmt.Sprintf(`{"status":%d, "friends":{"fans":%s, "focus" : %s}}`, common.STATUS_OK, pFans.StringFriendsArray(), pFocus.StringFriendsArray())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Get personal information.
//When parameter "msg" is 1, then the first top tweets will be
//sent to client.
func handlePersonal(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chUserInfo := make(chan *personal.UserInfo)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handlePersonal: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()
	args := r.URL.Query()
	uid := args[common.UID][0] //Get my-id
	fri := args[common.FRI][0] //An id of friend who will be checked.
	msg := args[common.MSG][0] //When "msg" is 1, then the first top tweets

	u, _ := strconv.Atoi(uid)
	f, _ := strconv.Atoi(fri)
	m, _ := strconv.Atoi(msg)

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	go personal.UserInformation(cxt, session, access_token, u, f, chUserInfo)
	pUserInfo := <-chUserInfo

	s := ""
	if m != 1 { //When "msg" is 1, then the first top tweets
		s = fmt.Sprintf(`{"status":%d, "user":%s}`, common.STATUS_OK, pUserInfo)
	} else {
		chTweetList := make(chan *tweet.TweetsList)
		go tweet.TweetList(cxt, f, session, access_token, 1, chTweetList)
		pTweetsList := <-chTweetList
		
		s = fmt.Sprintf(`{"status":%d, "user":%s, "tweets" : %s}`, common.STATUS_OK, pUserInfo, pTweetsList.StringTweetsArray()) 
	}
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Update relation between me and user.
func handleUpdateRelation(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chRes := make(chan *common.Result)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleUpdateRelation: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()
	args := r.URL.Query()
	fri := args[common.FRI][0] //An id of friend who will be checked.
	rel := args[common.REL][0] //Update relation between me and user:  0-cancle，1-focus

	re, _ := strconv.Atoi(rel)
	f, _ := strconv.Atoi(fri)

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	go personal.UpdateReleation(cxt,  session, access_token, f, re, chRes)
	pRes := <-chRes
	s := fmt.Sprintf(`{"status":%d, "result":%s}`, common.STATUS_OK, pRes.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}
