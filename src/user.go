package osc

import (
	"appengine"
	"appengine/urlfetch"

	"bytes"
	"fmt"

	"net/http"

	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"io"
	"io/ioutil"
)

type Token struct {
	UID          int    `json:"uid"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
}

func (self *Token) String() (s string) {
	json, _ := json.Marshal(self)
	s = string(json)
	return
}

type Logined struct {
	Cookie *http.Cookie
	Token  *Token
}

func (self *Logined) String() (s string) {
	s = self.Token.String()
	return
}

type OscUser struct {
	Account  string
	Password string
	AppId    string
	AppSec   string
}

func newOscUser(account, password, appId, appSec string) (usr *OscUser) {
	usr = new(OscUser)
	usr.Account = account
	usr.Password = password
	usr.AppId = appId
	usr.AppSec = appSec
	return
}

func (self *OscUser) buildLoginBody() (body string) {
	body = fmt.Sprintf(`email=%s&pwd=%s`, self.Account, self.Password)
	return
}

func (self *OscUser) buildOAuth2Body() (body string) {
	body = fmt.Sprintf(`client_id=%s&response_type=code&redirect_uri=%s&scope=%s&state=""&user_oauth_approval=true&email=%s&pwd=%s`, self.AppId, REDIRECT_URL, SCOPE, self.Account, self.Password)
	return
}

func (self *OscUser) login(cxt appengine.Context, ch chan *Logined) {
	defer func() {
		if e := recover(); e != nil {
			close(ch)
		}
	}()

	h := sha1.New()
	io.WriteString(h, self.Password)
	self.Password = hex.EncodeToString(h.Sum(nil))

	pClient := urlfetch.Client(cxt)
	body := self.buildLoginBody()

	pLogined := new(Logined)
	if r, e := http.NewRequest("POST", "https://www.oschina.net/action/user/hash_login", bytes.NewBufferString(body)); e == nil {

		r.Header.Add("Accept", "*/*")
		r.Header.Add("Accept-Encoding", ACCEPT_ENCODING)
		r.Header.Add("Accept-Language", ACCEPT_LANG)
		r.Header.Add("Connection", KEEP_ALIVE)
		r.Header.Add("Content-Type", API_REQTYPE)
		r.Header.Add("Host", OSC)
		r.Header.Add("Origin", ORIGINAL)
		r.Header.Add("User-Agent", AGENT)
		r.Header.Add("X-Requested-With", XMLHTTPREQUEST)
		r.Header.Add("Referer", AUTH_REF_URL)

		if resp, e := pClient.Do(r); e == nil {
			//Get cookie, and do OAuth2 in order to fetching "code".
			for _, v := range resp.Cookies() {
				if v.Value != "" {
					pLogined.Cookie = v
					code := self.oAuth2(pClient, v)
					pLogined.Token = self.getToken(pClient, code)
					break
				}
			}
		} else {
			panic(e)
		}
	} else {
		panic(e)
	}

	ch <- pLogined
}

func (self *OscUser) oAuth2(pClient *http.Client, cookie *http.Cookie) (code string) {
	body := self.buildOAuth2Body()

	if r, e := http.NewRequest("POST", AUTH_URL, bytes.NewBufferString(body)); e == nil {
		r.Header.Add("Accept", "*/*")
		r.Header.Add("Accept-Encoding", ACCEPT_ENCODING)
		r.Header.Add("Accept-Language", ACCEPT_LANG)
		r.Header.Add("Connection", KEEP_ALIVE)
		r.Header.Add("Content-Type", API_REQTYPE)
		r.Header.Add("Host", OSC)
		r.Header.Add("X-Requested-With", XMLHTTPREQUEST)
		r.Header.Add("User-Agent", AGENT)
		r.Header.Add("Referer", AUTH_REF_URL)
		r.Header.Add("Pragma", NO_CACHE)
		r.Header.Add("Cache-Control", NO_CACHE)
		r.Header.Add("Cache-Control", NO_CACHE)
		r.Header.Add("Cookie", "oscid="+cookie.Value)

		if resp, e := pClient.Do(r); e == nil {
			args := resp.Request.URL.Query()
			code = args["code"][0]
		} else {
			panic(e)
		}
	} else {
		panic(e)
	}
	return
}

func (self *OscUser) getToken(pClient *http.Client, code string) (pToken *Token) {
	body := fmt.Sprintf(TOKEN_BODY, APP_ID, APP_SEC, GRANT_TYPE, REDIRECT_URL, code, RET_TYPE)
	if r, e := http.NewRequest("POST", TOKEN_URL, bytes.NewBufferString(body)); e == nil {
		r.Header.Add("Content-Type", API_REQTYPE)
		if resp, e := pClient.Do(r); e == nil {
			pToken = new(Token)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e = json.Unmarshal(bytes, pToken); e != nil {
					pToken = nil
				}
			} else {
				panic(e)
			}
		} else {
			panic(e)
		}
	} else {
		panic(e)
	}
	return
}
