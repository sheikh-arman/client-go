package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"time"
)

var jwtkey = []byte("Neaj's Secret Key, He will not share it")
var TokenAuth *jwtauth.JWTAuth
var tokenString string
var token jwt.Token

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Item struct {
	Title string `json:"title"`
	Post  string `json:"post"`
	Id    int    `json:"id"`
}

var ID int
var feeds []Item
var Credslist map[string]string

func InitCred() {
	TokenAuth = jwtauth.New(string(jwa.HS256), jwtkey, nil)
	Credslist = make(map[string]string)
	creds := Credentials{
		"arman",
		"123",
	}
	Credslist[creds.Username] = creds.Password
}
func InitID() {
	ID = 0
}
func InitDB() {
	var feed Item
	feed = Item{
		Id:    ID,
		Title: "Nothing",
		Post:  "Lorem Ipsum Doller Site",
	}
	ID++
	feeds = append(feeds, feed)

	feed = Item{
		Id:    ID,
		Title: "Nothing2",
		Post:  "Lorem Ipsum Doller Site2",
	}
	ID++
	feeds = append(feeds, feed)

	feed = Item{
		Id:    ID,
		Title: "Nothing3",
		Post:  "Lorem Ipsum Doller Site3",
	}
	ID++
	feeds = append(feeds, feed)
}

func GetNewsFeeds(w http.ResponseWriter, r *http.Request) {
	log.Println("test")
	w.Header().Set("Content-Type", "application/json")
	sort.SliceStable(feeds, func(i, j int) bool {
		return feeds[i].Id < feeds[j].Id
	})
	json.NewEncoder(w).Encode(feeds)
}
func GetNewsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := chi.URLParam(r, "id")
	paramsID, _ := strconv.Atoi(param)

	for _, curFeed := range feeds {
		if curFeed.Id == paramsID {
			json.NewEncoder(w).Encode(curFeed)
		}
	}
}
func CreateNewsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newFeed Item
	err := json.NewDecoder(r.Body).Decode(&newFeed)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	newFeed.Id = ID
	feeds = append(feeds, newFeed)
	ID++
}
func DeleteNewsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := chi.URLParam(r, "id")
	paramsID, _ := strconv.Atoi(param)
	for index, curFeed := range feeds {
		if curFeed.Id == paramsID {
			feeds = append(feeds[:index], feeds[index+1:]...)
			break
		}
	}
	fmt.Println(feeds)
	json.NewEncoder(w).Encode(feeds)
}
func UpdateNewsFeed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	param := chi.URLParam(r, "id")
	paramID, _ := strconv.Atoi(param)
	var newFeed Item
	err := json.NewDecoder(r.Body).Decode(&newFeed)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
	for index, curFeed := range feeds {
		if curFeed.Id == paramID {
			feeds = append(feeds[:index], feeds[index+1:]...)
			newFeed.Id = paramID
			feeds = append(feeds, newFeed)
			sort.SliceStable(feeds, func(i, j int) bool {
				return feeds[i].Id < feeds[j].Id
			})
			json.NewEncoder(w).Encode(feeds)
			return
		}
	}
	json.NewEncoder(w).Encode("No data")
}
func Login(w http.ResponseWriter, r *http.Request) {
	var creds Credentials
	err := json.NewDecoder(r.Body).Decode(&creds)

	fmt.Println(creds)

	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	correctPassword, ok := Credslist[creds.Username]

	if !ok || creds.Password != correctPassword {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	expiretime := time.Now().Add(10 * time.Minute)

	_, tokenString, err := TokenAuth.Encode(map[string]interface{}{
		"aud": "arman",
		"exp": expiretime.Unix(),
	})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Value:   tokenString,
		Expires: expiretime,
	})
	w.WriteHeader(http.StatusOK)
}
func Logout(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "jwt",
		Expires: time.Now(),
	})
	w.WriteHeader(http.StatusOK)
}
func main() {
	InitCred()
	InitID()
	InitDB()

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/login", Login)

	r.Group(func(r chi.Router) {
		r.Use(jwtauth.Verifier(TokenAuth))
		r.Use(jwtauth.Authenticator)
		r.Route("newsfeeds", func(r chi.Router) {
			r.Get("/", GetNewsFeeds)
			r.Get("/{id}", GetNewsFeed)
			r.Post("/", CreateNewsFeed)
			r.Delete("/{id}", DeleteNewsFeed)
			r.Put("/{id}", UpdateNewsFeed)
		})
		r.Post("/logout", Logout)
	})
	port := 3000
	Server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}
	Server.ListenAndServe()
}
