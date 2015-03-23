package personal

import (
	"common"
	"encoding/json"
	"encoding/xml"

	"appengine"
	"appengine/urlfetch"

	"bytes"
	"fmt"

	"io/ioutil"
	"net/http"
)

func UserInformation(cxt appengine.Context, session string, uid int, friend int, ch chan *UserInfo) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.USER_INFORMATION_SCHEME, uid, friend, "")
	if r, e := http.NewRequest(common.POST, common.USER_INFORMATION_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pInfo := new(UserInfo)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := xml.Unmarshal(bytes, pInfo); e == nil {
					ch <- pInfo
				} else {
					ch <- nil
					cxt.Errorf("Error but still going: %v", e)
				}
			} else {
				ch <- nil
				panic(e)
			}
		} else {
			ch <- nil
			cxt.Errorf("Error but still going: %v", e)
		}
	} else {
		ch <- nil
		panic(e)
	}
}

func MyInformation(cxt appengine.Context, session string, uid int, ch chan *MyInfo) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.MY_INFORMATION_SCHEME, uid)
	if r, e := http.NewRequest(common.POST, common.MY_INFORMATION_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pMyInfo := new(MyInfo)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := xml.Unmarshal(bytes, pMyInfo); e == nil {
					ch <- pMyInfo
				} else {
					ch <- nil
					cxt.Errorf("Error but still going: %v", e)
				}
			} else {
				ch <- nil
				panic(e)
			}
		} else {
			ch <- nil
			cxt.Errorf("Error but still going: %v", e)
		}
	} else {
		ch <- nil
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
				if e := json.Unmarshal(bytes, pRes); e == nil {
					//pRes.Message =  string(bytes)
					ch <- pRes
				} else {
					ch <- nil
					panic(e)
				}
			} else {
				ch <- nil
				panic(e)
			}
		} else {
			ch <- nil
			panic(e)
		}
	} else {
		ch <- nil
		panic(e)
	}
}
