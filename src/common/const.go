package common

const (
	OSC             = "www.oschina.net"
	ORIGINAL        = "https://www.oschina.net"
	HOST            = ORIGINAL + "/action/openapi/"
	INT_HOST        = ORIGINAL + "/action/api/"
	API_REQTYPE     = "application/x-www-form-urlencoded; charset=UTF-8" //Request type
	API_RESTYPE     = "application/json"                                 //Response types
	POST            = "POST"
	GET             = "GET"
	KEEP_ALIVE      = "Keep-Alive"
	ACCEPT_LANG     = "zh-cn,zh;q=0.8,en-us;q=0.5,en;q=0.3"
	ACCEPT_ENCODING = "gzip,deflate,sdch"
	NO_CACHE        = "no-cache"
	XMLHTTPREQUEST  = "XMLHttpRequest"
	AGENT           = "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36"

	DATA_TYPE = "&dataType=json"

	TWEET_LIST_URL    = HOST + "tweet_list"
	TWEET_LIST_SCHEME = "user=%d&access_token=%s&pageIndex=%d&pageSize=25" + DATA_TYPE

	TWEET_PUB_URL    = HOST + "tweet_pub"
	TWEET_PUB_SCHEME = "access_token=%s&msg=%s" + DATA_TYPE

	COMMENT_REPLY_URL    = HOST + "comment_reply"
	COMMENT_REPLY_SCHEME = "id=%d&catalog=%d&content=%s&receiver=%d&authorid=%d&replyid=%d&isPostToMyZone=0&access_token=%s" + DATA_TYPE

	TWEET_DETAIL_URL    = HOST + "tweet_detail"
	TWEET_DETAIL_SCHEME = "id=%d&access_token=%s" + DATA_TYPE

	PERSONAL_FRIENDS_LIST_URL    = INT_HOST + "friends_list"
	PERSONAL_FRIENDS_LIST_SCHEME = "uid=%d&page=%d&relation=%d&pageSize=%d"

	USER_INFORMATION_URL    = INT_HOST + "user_information"
	USER_INFORMATION_SCHEME = "uid=%d&hisuid=%d&hisname=%s&pageIndex=0&pageSize=0"

	FOCUS_USER_URL    = HOST + "update_user_relation"
	FOCUS_USER_SCHEME = "friend=%d&relation=%d&access_token=%s" + DATA_TYPE

	MY_INFORMATION_URL    = INT_HOST + "my_information"
	MY_INFORMATION_SCHEME = "uid=%d"

	COMMENT_PUB_URL    = HOST + "comment_pub"
	COMMENT_PUB_SCHEME = "catalog=%d&id=%d&content=%s&access_token=%s" + DATA_TYPE

	COMMENT_LIST_URL    = HOST + "comment_list"
	COMMENT_LIST_SCHEME = "catalog=%d&id=%d&page=%d&pageSize=99&access_token=%s" + DATA_TYPE

	ACTIVE_LIST_URL    = HOST + "active_list"
	ACTIVE_LIST_SCHEME = "catalog=%d&user=%d&page=%d&pageSize=99&access_token=%s" + DATA_TYPE

	CLEAR_NOTICE_URL    = HOST + "clear_notice"
	CLEAR_NOTICE_SCHEME = "type=%d&access_token=%s" + DATA_TYPE

	API_USER_URL    = HOST + "user"
	API_USER_SCHEME = "access_token=%s" + DATA_TYPE

	ID    = "id"
	UID   = "uid"
	IDENT = "ident"
	PAGE  = "page"
	MSG   = "msg"
	FRI   = "fri"
	REL   = "rel"
	ME    = "me"
	MP    = "mp"
	FP    = "fp"

	STATUS_OK  = 200
	STATUS_ERR = 300

	SUCCESS    = "1"  //success
	DUPLICATED = "-1" //duplicate
	NO_LOGIN   = "0"  //No login

	LOGIN_URL    = ORIGINAL + "/action/user/hash_login"
	AUTH_URL     = ORIGINAL + "/action/oauth2/authorize"
	TOKEN_URL    = HOST + "token"
	TOKEN_BODY   = "client_id=%s&client_secret=%s&grant_type=%s&redirect_uri=%s&code=%s&dataType=%s"
	AUTH_REF_URL = ORIGINAL + "/action/oauth2/authorize?response_type=code&client_id=%s&redirect_uri=%s"
	GRANT_TYPE   = "authorization_code"
	RET_TYPE     = "json"
	EDIT_URL     = "/admin/profile"
)
