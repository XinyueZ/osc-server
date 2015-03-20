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

		tweetsMp := make(map[int]*tweet.Tweet)
		ch2 := make(chan *tweet.Tweet, l)

		for _, elem := range tweetFavorites {
			go tweet.TweetDetail(cxt, session, access_token, elem.ObjectId, ch2)
		}

		mount := 0
		for index := 0; index < l; index++ {
			item := <-ch2
			if item != nil {
				tweetsMp[item.Id] = item
				mount++
			}
		}

		ret := make([]*tweet.Tweet, mount)
		pos := 0 //Loop to mount
		for index := 0; index < l; index++ {
			elem := tweetFavorites[index]
			item := tweetsMp[elem.ObjectId]
			if item != nil {
				ret[pos] = item
				pos++
			}
		}

		json, _ := json.Marshal(&ret)
		ch1 <- string(json)
	} else {
		ch1 <- "null"
	}
}
