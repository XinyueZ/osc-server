package osc

/*
* Define application client-id(APP_ID) and client-security(APP_SEC) code here which will be gotten when you register an application at www.oschina.net
* Also the key and commonIV for encryption, they must be equal to the same file in client code, the "sec.java".
 */
const (
	APP_ID  = "APP_ID"
	APP_SEC = "APP_SEC"
)

var KEY []byte = []byte{11, 45}
var IV []byte = []byte{35, 23}
