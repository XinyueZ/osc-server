package favorite

import (
	"appengine"
	"appengine/datastore"
	"strconv"
	"time"
)

func AddTweetFavorite(cxt appengine.Context, session string, access_token string, id int, objectId int) {
	editTime, _ := strconv.ParseInt(time.Now().Local().Format("20060102150405"), 10, 64)
	pTweetFavorite := &TweetFavorite{id, objectId, editTime}
	datastore.Put(cxt, datastore.NewIncompleteKey(cxt, "TweetFavorite", nil), pTweetFavorite)
}
