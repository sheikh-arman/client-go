package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth"
	"github.com/lestrrat-go/jwx/jwa"
	"github.com/lestrrat-go/jwx/jwt"
	"net/http"
	"strconv"
)

var jwtkey = []byte("Neaj's Secret Key, He will not share it")
var TokenAuth *jwtauth.JWTAuth
var tokenString string
var token jwt.Token

type Credentials struct{
	Username string `json:"username"`
	Password  string `json:"password"`
}



type Item struct {
	Title string `json:"title"`
	Post  string `json:"post"`
	Id   int `json:"id"`
}

var ID int
var feeds []Item
var Credslist map[string]string

func InitCred(){
	TokenAuth = jwtauth.New(string(jwa.HS256), jwtkey, nil)
	Credslist = make(map[string]string)
	creds := Credentials{
		"arman",
		"123",
	}
	Credslist[creds.Username]=creds.Password
}
func InitID(){
	ID=0
}
func InitDB(){
	var feed Item
	feed=Item{
		Id: ID,
		Title: "Nothing"
		Post: "Lorem Ipsum Doller Site"
	}
	ID++
	feeds=append(feeds,feed)

	feed=Item{
		Id: ID,
		Title: "Nothing2"
		Post: "Lorem Ipsum Doller Site2"
	}
	ID++
	feeds=append(feeds,feed)

	feed=Item{
		Id: ID,
		Title: "Nothing3"
		Post: "Lorem Ipsum Doller Site3"
	}
	ID++
	feeds=append(feeds,feed)
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
		r.Use(jwtauth.Verifier(data.TokenAuth))
		r.Use(jwtauth.Authentication)
		r.Router("newsfeeds", func(r chi.Router) {
			r.Get("/", GetNewsfeeds)
			r.Get("/{id}", GetNewsFeed)
			r.Post("/", CreateNewsFeed)
			r.Delete("/{id}", DeleteNewsFeed)
			r.Put("/{id}", UpdateNewsFeed)
		})
		r.Post("/logout", Logout)
	})

	Server := &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: r,
	}
	fmt.Println(Server.ListenAndServe())
}
