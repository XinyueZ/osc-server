package favorite

import (
	"appengine"
	"appengine/datastore"
	"encoding/json"
	"sort"
  "tweet"
)

func GetTweetFavoritesList(cxt appengine.Context, session string, access_token string, id int, ch1 chan string) {
	q := datastore.NewQuery("TweetFavorite").Filter("Id =", id)
	tweetFavorites := make(TweetFavorites, 0)
	keys, _ := q.GetAll(cxt, &tweetFavorites)

	l := len(keys)
	if l > 0 {
    sort.Sort(tweetFavorites)

		tweets := make(map[int]*tweet.Tweet)
		ch2    := make(chan *tweet.Tweet, l)

		for _, elem := range tweetFavorites {
      go tweet.TweetDetail (cxt, session, access_token, elem.ObjectId, ch2 )
		}

		for index:= 0; index < l; index++ {
			item := <-ch2
			tweets[item.Id] = item
		}


		ret :=  make([]*tweet.Tweet, l)
		for index:= 0; index < l; index++ {
			elem := tweetFavorites[index]
			ret[index] = tweets[elem.ObjectId]
		}

		json, _ := json.Marshal(&ret)
		ch1 <- string(json)
	} else {
		ch1 <- "null"
	}
}
