package osc

const (
	OSC         = "www.oschina.net"
	HOST        = "http://www.oschina.net/action/api/"
	API_REQTYPE = "application/x-www-form-urlencoded" //Request type
	API_RESTYPE = "application/json"                  //Response types
	POST        = "POST"
	GET         = "GET"
	KEEP_ALIVE  = "Keep-Alive"

	LOGIN_SCHEME        = `username=%s&pwd=%s&keep_login=1`
	LOGIN_VALIDATE_HTTP = HOST + "login_validate"

	TWEET_LIST        = HOST + "tweet_list?uid=%d&pageIndex=%d&pageSize=25"
	TWEET_LIST_SCHEME = "uid=%d&pageIndex=%d&pageSize=25"
	TWEET_PUB         = HOST + "tweet_pub"
	TWEET_PUB_SCHEME  = "uid=%d&msg=%s"

	ACCOUNT    = "u"
	PWD        = "pw"
	UID        = "uid"
	PAGE       = "page"
	MSG        = "msg"
	STATUS_OK  = 200
	STATUS_ERR = 300

	SUCCESS    = "1"  //success
	DUPLICATED = "-1" //duplicate
	NO_LOGIN   = "0"  //No login
)
