package handler

import (
	"Api-Server/platform/newsfeed"
	"encoding/json"
	"net/http"
)

func NewsfeedGet(feed *newsfeed.Repo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		items := feed.GetAll()
		json.NewEncoder(w).Encode(items)
	}
}
