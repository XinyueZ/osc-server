package personal

import (
	"bytes"
	"common"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"appengine"
	"appengine/urlfetch"
)

//The API-User information.
type Me struct {
	Id       int    `json:"id"`
	Location string `json:"location"`
	Email    string `json:"email"`
	Name     string `json:"name"`
	Gender   string `json:"gender"`
	Avatar   string `json:"avatar"`
	Url      string `json:"url"`
}

func (self Me) String() (s string) {
	json, _ := json.Marshal(&self)
	s = string(json)
	return
}

func GetMe(cxt appengine.Context, session string, access_token string, ch chan *Me) {
	client := urlfetch.Client(cxt)
	body := fmt.Sprintf(common.API_USER_SCHEME, access_token)
	if r, e := http.NewRequest(common.POST, common.API_USER_URL, bytes.NewBufferString(body)); e == nil {
		common.MakeHeader(r, "oscid="+session, 0)
		if resp, e := client.Do(r); e == nil {
			if resp != nil {
				defer resp.Body.Close()
			}
			pMe := new(Me)
			if bytes, e := ioutil.ReadAll(resp.Body); e == nil {
				if e := json.Unmarshal(bytes, pMe); e == nil {
					ch <- pMe
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
