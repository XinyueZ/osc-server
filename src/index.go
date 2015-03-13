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

	"comment"
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
	http.HandleFunc("/tweetReply", handleTweetReply)
	http.HandleFunc("/tweetCommentPub", handleTweetCommentPub)
	http.HandleFunc("/tweetCommentList", handleTweetCommentList)
	http.HandleFunc("/atMeNoticesList", handleAtMeNoticesList)
	http.HandleFunc("/tweetDetail", handleTweetDetail)
	http.HandleFunc("/friendsList", handleFriendsList)
	http.HandleFunc("/userInformation", handlePersonal)
	http.HandleFunc("/updateRelation", handleUpdateRelation)
	http.HandleFunc("/myInformation", handleMyInformation)
	http.HandleFunc("/newCommentsNoticesList", handleNewCommentsNoticesList)
	http.HandleFunc("/clearAtNotice", handleClearAtNotice)
	http.HandleFunc("/clearCommentsNotice", handleClearCommentsNotice)

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
	appId := strings.TrimSpace((strings.Split(data[2], "="))[1])
	appSec := strings.TrimSpace((strings.Split(data[3], "="))[1])
	redirectUrl := strings.TrimSpace((strings.Split(data[4], "="))[1])
	scope := strings.TrimSpace((strings.Split(data[5], "="))[1])

	pUser := user.NewOscUser(account, pwd, appId, appSec, redirectUrl, scope)
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

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token
	msg := cookies[2].Value          //Message

	go tweet.TweetPub(cxt, session, access_token, msg, chTweetPub)
	pRes := <-chTweetPub

	s := fmt.Sprintf(`{"status":%d, "result":%s}`, common.STATUS_OK, pRes.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Reply someone of tweets.
func handleTweetReply(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweetReply := make(chan *common.Result)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetReply: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	cookies := r.Cookies()                          //Session in cookies passt
	session := cookies[0].Value                     //Get user-session
	access_token := cookies[1].Value                //Get user-token
	id, _ := strconv.Atoi(cookies[2].Value)         //The id of original object that contains comments , replys, ie. Tweet, Blog, etc.
	content := cookies[3].Value                     //The content to reply.
	receiverId, _ := strconv.Atoi(cookies[4].Value) //The userId of a person who will get my reply.
	authorId, _ := strconv.Atoi(cookies[5].Value)   //The userId of author who wirtes the reply.
	replyId, _ := strconv.Atoi(cookies[6].Value)    //The replied object, ie. a comment.

	pRes := comment.TweetCommentReply(cxt, session, access_token, id, content, receiverId, authorId, replyId, chTweetReply)

	s := fmt.Sprintf(`{"status":%d, "result":%s}`, common.STATUS_OK, pRes.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Publish comment to tweet.
func handleTweetCommentPub(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chResult := make(chan *common.Result)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetCommentPub: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	id := args[common.ID][0] //Get my-id

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token
	content := cookies[2].Value      //Comment- content

	i, _ := strconv.Atoi(id)

	go comment.TweetCommentPub(cxt, session, access_token, i, content, chResult)
	pRes := <-chResult

	s := fmt.Sprintf(`{"status":%d, "result":%s}`, common.STATUS_OK, pRes.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Get all friends.
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

	go personal.UpdateReleation(cxt, session, access_token, f, re, chRes)
	pRes := <-chRes
	s := fmt.Sprintf(`{"status":%d, "result":%s}`, common.STATUS_OK, pRes.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Get my personal information.
func handleMyInformation(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chMyInfo := make(chan *personal.MyInfo)
	chTweetActivesList := make(chan *personal.ActivesList)
	chCommentActivesList := make(chan *personal.ActivesList)
	chApiUser := make(chan *personal.Me)

	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleMyInformation: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	me, _:= strconv.Atoi(args[common.ME][0] ) //Show me? 0 not, !=0 yes.
	showMe := false
	if me > 0 {
		showMe = true
	}

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	go personal.MyInformation(cxt, session, access_token, chMyInfo)
	pMyInfo := <-chMyInfo

	//Get first page of active-list of at you of tweets that I joined.
	pTweetActivesList := personal.LastTweetActiveList(cxt, session, access_token, pMyInfo.Uid, 1, showMe, chTweetActivesList)
	sTweetActivesList := "null"
	if pTweetActivesList != nil {
		sTweetActivesList = pTweetActivesList.StringActivesArray()
	}
	//Get first page of active-list of comments of tweets that I've written.
	pCommentActivesList := personal.LastCommentActiveList(cxt, session, access_token, pMyInfo.Uid, 1, showMe, chCommentActivesList)
	sCommentActivesList := "null"
	if pCommentActivesList != nil {
		sCommentActivesList = pCommentActivesList.StringActivesArray()
	}

	go personal.GetMe(cxt, session, access_token, chApiUser)
	pMe := <-chApiUser
	homeUrl := ""
	editUrl := ""
	if pMe != nil {
		homeUrl = pMe.Url
		editUrl = pMe.Url + common.EDIT_URL
	}
	s := fmt.Sprintf(`{"status":%d, "am":%s,  "url" : {"home":"%s", "edit" : "%s"},  "atMe" : %s,  "comments":%s}`, common.STATUS_OK, pMyInfo, homeUrl, editUrl, sTweetActivesList, sCommentActivesList)
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Get all comments of a tweet.
func handleTweetCommentList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chCommentList := make(chan *comment.CommentList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetCommentList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	id := args[common.ID][0]     //Which tweet item
	page := args[common.PAGE][0] //Which page

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	i, _ := strconv.Atoi(id)
	pg, _ := strconv.Atoi(page)

	go comment.TweetCommentList(cxt, session, access_token, i, pg, chCommentList)
	pCommentList := <-chCommentList
	s := fmt.Sprintf(`{"status":%d, "comments":%s}`, common.STATUS_OK, pCommentList.StringCommentArray())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Show  tweet "at me" notices list.
func handleAtMeNoticesList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chActivesList := make(chan *personal.ActivesList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleAtMeNoticesList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	uid := args[common.UID][0]   //User id
	page := args[common.PAGE][0] //Which page
	me, _:= strconv.Atoi(args[common.ME][0] ) //Show me? 0 not, !=0 yes.
	showMe := false
	if me > 0 {
		showMe = true
	}

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	user, _ := strconv.Atoi(uid)
	pg, _ := strconv.Atoi(page)
	go personal.TweetActiveList(cxt, session, access_token, user, pg, showMe, chActivesList)
	pActivesList := <-chActivesList

	s := "null"
	if pActivesList != nil {
		s = pActivesList.StringActivesArray()
	}
	fmt.Sprintf(`{"status":%d, "notices":%s}`, common.STATUS_OK, s)
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Show new comment notices list.
func handleNewCommentsNoticesList(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chActivesList := make(chan *personal.ActivesList)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleNewCommentsNoticesList: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	uid := args[common.UID][0]   //User id
	page := args[common.PAGE][0] //Which page
	me, _:= strconv.Atoi(args[common.ME][0] ) //Show me? 0 not, !=0 yes.
	showMe := false
	if me > 0 {
		showMe = true
	}

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	user, _ := strconv.Atoi(uid)
	pg, _ := strconv.Atoi(page)
	go personal.CommentsActiveList(cxt, session, access_token, user, pg, showMe, chActivesList)
	pActivesList := <-chActivesList

	s := "null"
	if pActivesList != nil {
		s = pActivesList.StringActivesArray()
	}
	fmt.Sprintf(`{"status":%d, "notices":%s}`, common.STATUS_OK, s)
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Get detail of a single tweet.
func handleTweetDetail(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chTweet := make(chan *tweet.Tweet)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetDetail: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	args := r.URL.Query()
	id := args[common.ID][0] //Tweet id

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	i, _ := strconv.Atoi(id)

	go tweet.TweetDetail(cxt, session, access_token, i, chTweet)
	pTweet := <-chTweet

	s := fmt.Sprintf(`{"status":%d, "tweet":%s}`, common.STATUS_OK, pTweet.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Clear "@me" notice
func handleClearAtNotice(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chResult := make(chan *common.Result)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleTweetDetail: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	pRes := common.ClearAtNotice(cxt, session, access_token, chResult)
	s := fmt.Sprintf(`{"status":%d, "result":%s}`, common.STATUS_OK, pRes.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}

//Clear notice of comments.
func handleClearCommentsNotice(w http.ResponseWriter, r *http.Request) {
	cxt := appengine.NewContext(r)
	chResult := make(chan *common.Result)
	defer func() {
		if err := recover(); err != nil {
			cxt.Errorf("handleClearCommentsNotice: %v", err)
			fmt.Fprintf(w, `{"status":%d}`, common.STATUS_ERR)
		}
	}()

	cookies := r.Cookies()           //Session in cookies passt
	session := cookies[0].Value      //Get user-session
	access_token := cookies[1].Value //Get user-token

	pRes := common.ClearCommentsNotice(cxt, session, access_token, chResult)
	s := fmt.Sprintf(`{"status":%d, "result":%s}`, common.STATUS_OK, pRes.String())
	w.Header().Set("Content-Type", common.API_RESTYPE)
	fmt.Fprintf(w, s)
}
