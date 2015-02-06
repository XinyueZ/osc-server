package osc
 
const (
  OSC                 = "www.oschina.net"
  HOST                = "http://www.oschina.net/action/api/"
  API_REQTYPE         = "application/x-www-form-urlencoded" //Request type
  API_RESTYPE         = "application/json" //Response types
  POST                = "POST"
  GET                 = "GET"
  KEEP_ALIVE          = "Keep-Alive"
  LOGIN_SCHEME        = `username=%s&pwd=%s&keep_login=1`
  LOGIN_VALIDATE_HTTP = HOST + "login_validate"
  TWEET_LIST          = HOST + "tweet_list?uid=%s&pageIndex=%d&pageSize=25"
  TWEET_PUB           = HOST + "tweet_pub"
  TWEET_PUB_SCHEME    = "uid=%s&msg=%s"
)
