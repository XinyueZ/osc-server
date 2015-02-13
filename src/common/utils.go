package common


import (
  "fmt"
  "net/http"
  "net/url"
  "strconv"
)

// UrlEncoded encodes a string like Javascript's encodeURIComponent()
func UrlEncoded(str string) (string, error) {
    u, err := url.Parse(str)
    if err != nil {
        return "", err
    }
    return u.String(), nil
}




//Make header for some requests.
func MakeHeader(r *http.Request, cookie string, length int) {
  r.Header.Add("Content-Type", API_REQTYPE)
  r.Header.Add("Content-Length", strconv.Itoa(length))
  r.Header.Add("Host", OSC)
  r.Header.Add("Connection", KEEP_ALIVE)
  r.Header.Add("Cookie", cookie)
}

//Util method to print request header.
func PrintHeader(r *http.Request, w http.ResponseWriter) {
  header := r.Header
  for k, v := range header {
    fmt.Fprintf(w, "k:%s v:%s", k, v)
  }
}
