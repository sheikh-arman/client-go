package main

import (
	"Api-Server/handler"
	"Api-Server/platform/newsfeed"
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
)

func main() {
	port := ":3000"
	feed := newsfeed.New()
	feed.Add(newsfeed.Item{
		Title: "Hello",
		Post:  "World",
	})

	r := chi.NewRouter()

	r.Get("/newsfeed", handler.NewsfeedGet(feed))
	r.Post("/newsfeed", handler.NewsfeedPost(feed))

	fmt.Println("Serving on " + port)
	http.ListenAndServe(port, r)

}
