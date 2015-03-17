package favorite

//The API-User information.
type TweetFavorite struct {
	Id       int
	ObjectId int
	AddTime  int64
}

type TweetFavorites []TweetFavorite

func (s TweetFavorites) Len() int {
    return len(s)
}
func (s TweetFavorites) Swap(i, j int) {
    s[i], s[j] = s[j], s[i]
}
func (s TweetFavorites) Less(i, j int) bool {
    return  s[i].AddTime > s[j].AddTime 
}
