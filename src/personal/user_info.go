package personal

import (
	"common"
	"encoding/json"

	"appengine"
	"appengine/urlfetch"

	"bytes"
	"fmt"

	"io/ioutil"
	"net/http"
)

func UserInformation(cxt appengine.Context, session string, access_token string, uid int, friend int, ch chan *UserInfo) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.USER_INFORMATION_SCHEME, uid, friend, access_token)
	//fmt.Fprintf(w, `%s\n`, body)
	if r, e := http.NewRequest(common.POST, common.USER_INFORMATION_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		//fmt.Fprintf(w, `oscid=%s\n`, session)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pUserInfo := new(UserInfo)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				//fmt.Fprintf(w, `%s\n`, string(bytes))
				if err := json.Unmarshal(bytes, pUserInfo); err == nil {
					ch <- pUserInfo
				} else {
					panic(e)
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
}

//Update relation between me and friend.
//0-cancle，1-focus
func UpdateReleation(cxt appengine.Context, session string, access_token string, friend int, relation int, ch chan *common.Result) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.FOCUS_USER_SCHEME, friend, relation, access_token)
	//fmt.Fprintf(w, `%s\n`, body)
	if r, e := http.NewRequest(common.POST, common.FOCUS_USER_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		//fmt.Fprintf(w, `oscid=%s\n`, session)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pRes := new(common.Result)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				//fmt.Fprintf(w, `%s\n`, string(bytes))
				if err := json.Unmarshal(bytes, pRes); err == nil {
					//pRes.Message =  string(bytes)
					ch <- pRes
				} else {
					//pRes.Message =  string(bytes)
					ch <- pRes
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
}
