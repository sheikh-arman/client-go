package main

import (
	"Api-Server/handler"
	"Api-Server/platform/newsfeed"
	"encoding/json"
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

	r.Post("/newsfeed", func(w http.ResponseWriter, r *http.Request) {
		request := map[string]string{}
		json.NewDecoder(r.Body).Decode(request)

		feed.Add(newsfeed.Item{
			Title: request["title"],
			Post:  request["post"],
		})

		w.Write([]byte("Good job!"))
	})
	fmt.Println("Serving on " + port)
	http.ListenAndServe(port, r)

}
