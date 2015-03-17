package favorite

import (
	"appengine"
	"appengine/datastore"
)

func DelTweetFavorite(cxt appengine.Context, session string, access_token string, id int, objectId int) {
	q := datastore.NewQuery("TweetFavorite").Filter("Id=", id).Filter("ObjectId=", objectId)
	tweetFavorites := make([]TweetFavorite, 0)
	keys, _ := q.GetAll(cxt, &tweetFavorites)
	datastore.DeleteMulti(cxt, keys)
}
