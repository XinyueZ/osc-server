osc-server API
==============

Encapsulation for the [openAPI](http://http://www.oschina.net/openapi) of  [www.oschina.net](http://www.oschina.net) 

Host on Google AppEngine.

Written in [Golang](http://www.golang.org).

##Why this API

The problem of original [openAPI](http://http://www.oschina.net/openapi)  is that users have to finish authentication by using webview or browser and be forced to redirect to web-app to be granted access-token, the web-app feeling. If they develop a mobile client only, no chance, they have to be redirected. That's the OAuth2 model of [oschina](http://www.oschina.net).  

There isn't any client library to ease users' authentication, either.

SAY NO to a hybrid-app.

##Profit 

- Bypass the WebView or Browser of OAuth2.
- More customized programming interfaces for client applications.
- An in dev-progress [Java-Lib](https://github.com/XinyueZ/osc-tweet/tree/master/osc4j) of this encapsulation is available.
- High speed feeling for oversee Chinese users.

##Host on  AppEngine

- Checkout whole repository.
- Register an application [here](http://www.oschina.net/openapi/client/edit). Don't worry about Chinese language, the Google's translator fills your willing.
- You got an application-ID and an application-private-key after registration. 
- Find [sec.go](https://github.com/XinyueZ/osc-server/blob/master/src/common/sec.go) and input APP_ID for the application-ID, APP_SEC for the application-private-key.
- Input KEY and IV in [sec.go](https://github.com/XinyueZ/osc-server/blob/master/src/common/sec.go) for login encryption. See below section for login.


Update   |  API| Method
--------|---------|---------
02-26-2015             |  [/userInformation](#user-informationuser_information)|GET  
  | [/updateRelation](#tweettweet_list)|POST 
02-20-2015             |  [/friendsList](#friends-listfriends_list)|GET  
02-13-2015|  [/login](#login)|POST
  | [/tweetList](#tweettweet_list)|GET 
  | [/myTweetList](#tweettweet_list)|GET 
  | [/hotspotTweetList](#tweettweet_list) |GET 
 
##Common

Var   |  Value
--------|---------
Host             |osc-server-848.appspot.com
 Request type|application/x-www-form-urlencoded; charset=UTF-8
 
##Login
API: POST  /login

Request body:

Var     | Type    | Comment
--------|---------|---------
u       |string   |login account
pw      |string   |password

Example:
Body of plain-text in request

```java
u=helloworld_account&pw=4567789 
```

######The body must be encrypted with AES/CFB/Nopadding mode.The key and commonIV must be equal to the values in [sec.go](https://github.com/XinyueZ/osc-server/blob/master/src/common/sec.go).

Return feeds:

Var      | Type     | Comment
---------|---------|---------
status   |int     |See [Status code](#status-code)
user   |struct     |A user.

Struct of a user.

Var      | Type     | Comment
---------|---------|---------
uid        |int   |An user id of [oschina](http://www.oschina.net) internal.
expired   |int   |Time to expire current session in seconds.

Example:

```json
{
	"status": 200,
	"user": {
		"uid": 113101,
		"expired": 517453
	}
}
```
Return cookies:

Var   | Type       | Comment
--------|---------|---------
oscid              |string  |Session Id.
access_token              |string  |Access-Token for current user.
Example:
```
oscid=asdfasdfw5w456esgsdfg&pw=23434-456657dfg-ezt457-ert 
```

##Tweet(tweet_list)

Get list of tweets.

APIs:

Var   | Method   |  Comment
--------| --------|---------
/tweetList   | GET    | Get all tweets.
/myTweetList   | GET    | Get only my tweets.
/hotspotTweetList   | GET    | Get hotspot tweets.

Request parameters:

Var   | Type       | Available API| Comment
--------|---------|---------|---------
page             |int  |All APIs|Page number, start at 1.
uid              |int  |/myTweetList|An user id of [oschina](http://www.oschina.net) internal.

Request cookie:

Var   | Type       | Comment
--------|---------|---------
oscid              |string  |Session Id after login.
access_token              |string  |Access-Token for current user after login.

Example:
```
oscid=asdfasdfw5w456esgsdfg&pw=23434-456657dfg-ezt457-ert 
```
Return:

Var      | Type     | Comment
---------|---------|---------
status   |int     |See [Status code](#status-code)
tweets   |array of struct     |All tweets.

Struct of a tweet.

Var      | Type     | Comment
---------|---------|---------
id        |int   |Message Id.
portrait        |string   |Author photo thumbnail url.
author        |string   |Author name.
authorid        |int   |Author Id in [oschina](http://www.oschina.net) internal.
body        |string   |Rich text of the tweet message
commentCount        |int   |Count of all comments on the topic.
pubDate        |string   |Time of publishing this tweet.
imgSmall        |string   |Url to the small image attachment of this tweet.
imgBig        |string   |Url to the big image attachment of this tweet.

Example:
```json
{
	"status": 200,
	"tweets": [{
		"id": 4877440,
		"portrait": "http://static.oschina.net/uploads/user/48/97321_50.jpg",
		"author": "MusterMann",
		"authorid": 97321,
		"body": "asfasdfasdfasdfasdfsdfsdaf",
		"commentCount": 0,
		"pubDate": "2015-02-13 17:07:53",
		"imgSmall": "",
		"imgBig": ""
	}, {
		"id": 4877437,
		"portrait": "http://static.oschina.net/uploads/user/51/102723_50.jpg?t=1411184780000",
		"author": "zhuxinyu",
		"authorid": 102723,
		"body": "9dui bu qi, dong dan , ta de hao wo dao le9",
		"commentCount": 0,
		"pubDate": "2015-02-13 17:07:11",
		"imgSmall": "",
		"imgBig": ""	
  },
  .................
  ]
}
```

##Friends-List(friends_list)

A list of my friends(users on [oschina](http://www.oschina.net) ), including who focus on me and my fans.

API: GET  /friendsList

Request cookie:

Var   | Type       | Comment
--------|---------|---------
oscid              |string  |Session Id after login.
access_token              |string  |Access-Token for current user after login.

Example:
```
oscid=asdfasdfw5w456esgsdfg&pw=23434-456657dfg-ezt457-ert 
```
Return:

Var      | Type     | Comment
---------|---------|---------
status   |int     |See [Status code](#status-code)
friends        |struct   |Struct of fans and focus
fans        |array of struct   |Fans of me.
focus        |array of struct   |Who have focused on me.

Struct of friend:

Var      | Type     | Comment
---------|---------|---------
expertise        |string   |User skills: Android developer, iOS developer etc.
name        |string   |User name.
userid        |int   |User id of [oschina](http://www.oschina.net) internal.
gender        |int   |1:Male 2:Female
portrait        |string   |Author photo thumbnail url.

Example:

```json
{
	"status": 200,
	"friends": {
		"fans": [{
			"expertise": "asdfasfasdf",
			"name": "ertwertwetwert",
			"userid": 345,
			"gender": 1,
			"portrait": "http://static.oschina.net/uploads/user/428/wertwetwertwert.jpg?t=1388468572000"
		}, {
			"expertise": "dagfertwrtwetertwert",
			"name": "wertwetwetwert",
			"userid": 345345,
			"gender": 2,
			"portrait": "http://www.oschina.net/img/portrait.gif"
		}, 
		.........
		],
		"focus": [{
			"expertise": "bsdfgfdgsdfg",
			"name": "wertwetwetwetwetwert",
			"userid": 356456,
			"gender": 1,
			"portrait": "http://www.oschina.net/img/portrait.gif"
		},
		.........
		]
	}
}
```

##User-Information(user_information)

Get personal information of a user.

API: GET  /userInformation

Request cookie:

Var   | Type       | Comment
--------|---------|---------
oscid              |string  |Session Id after login.
access_token              |string  |Access-Token for current user after login.

Example:
```
oscid=asdfasdfw5w456esgsdfg&pw=23434-456657dfg-ezt457-ert 
```

Request parameters:

Var   | Type       |  Comment
--------|---------|---------
uid             |int  |My user-id of [oschina](http://www.oschina.net) internal.
fri              |int  |An user-id of [oschina](http://www.oschina.net) internal, whose information will be checked.
msg             |int  |1: Show last top tweets, might be 20 lines. non-1: no tweets.


Return:

Var      | Type     | Comment
---------|---------|---------
status   |int     |See [Status code](#status-code)
uid        |int   |User id of [oschina](http://www.oschina.net) internal.
name        |string   |User-name.
ident        |string   |User-nickname.
province        |string   |From province
city        |string   |From city.
platforms        |string   |Favorite platforms to work.
expertise        |string   |Skill.
portrait        |string   |Author photo thumbnail url.
gender        |int   |1:Male 2:Female
relation        |int   |1- has been focuse 2-focused each other 3-no any relation.
tweets | array of struct   |See [/tweetList](#tweettweet_list)

Example:

```json
{
	"status": 200,
	"user": {
		"uid": 3455,
		"name": "sdf",
		"ident": "sdf",
		"province": "Hongkou",
		"city": "Shanghai",
		"platforms": "Java EE Java SE Android HTML/CSS Linux/Unix ",
		"expertise": "Java",
		"portrait": "http://static.oschina.net/uploads/user/585/117033991_50.jpg?t=1408671734000",
		"gender": 1,
		"relation": 3
	},
	"tweets": [{
		"id": 4877440,
		"portrait": "http://static.oschina.net/uploads/user/48/97321_50.jpg",
		"author": "MusterMann",
		"authorid": 97321,
		"body": "asfasdfasdfasdfasdfsdfsdaf",
		"commentCount": 0,
		"pubDate": "2015-02-13 17:07:53",
		"imgSmall": "",
		"imgBig": ""
	}, {
		"id": 4877437,
		"portrait": "http://static.oschina.net/uploads/user/51/102723_50.jpg?t=1411184780000",
		"author": "zhuxinyu",
		"authorid": 102723,
		"body": "9dui bu qi, dong dan , ta de hao wo dao le9",
		"commentCount": 0,
		"pubDate": "2015-02-13 17:07:11",
		"imgSmall": "",
		"imgBig": ""	
  },
  .................
  ]
	
}
```

##Status code

Var   |   Comment
--------| ---------
200              | Success
300              | Fail

LICENSE
-

> The MIT License (MIT)
> 
> Copyright (c) 2015 Chris Xinyue Zhao
> 
> Permission is hereby granted, free of charge, to any person obtaining
> a copy of this software and associated documentation files (the
> "Software"), to deal in the Software without restriction, including
> without limitation the rights to use, copy, modify, merge, publish,
> distribute, sublicense, and/or sell copies of the Software, and to
> permit persons to whom the Software is furnished to do so, subject to
> the following conditions:
> 
> The above copyright notice and this permission notice shall be
> included in all copies or substantial portions of the Software.
> 
> THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
> EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
> MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
> IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY
> CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT,
> TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE
> SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

