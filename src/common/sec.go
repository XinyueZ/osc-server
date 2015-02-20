package common

/*
* Define application client-id(APP_ID) and client-security(APP_SEC) code here which will be gotten when you register an application at www.oschina.net
* Also the key and commonIV for encryption, they must be equal to the same file in client code, the "sec.java".
 */
const (
	APP_ID       = "EOW46fNRtr7FgSlHAVz4"
	APP_SEC      = "eXc7xJv1uliOm3WRmo7r9IwzuqrvxYYu"
	REDIRECT_URL = "http://www.oschina.net"
)

var KEY []byte = []byte{11, 45, 78, 110, 118, 9, 3, 4, 18, 47, 3, 7, 77, 8, 56, 101}
var IV []byte = []byte{34, 35, 35, 57, 68, 4, 35, 36, 7, 8, 35, 23, 35, 86, 35, 23}
