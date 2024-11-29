package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

type router struct {
}

func (rt *router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/a":
		fmt.Fprintln(w, "executing /a")
	case "/b":
		fmt.Fprintln(w, "executing /b")
	case "/c":
		fmt.Fprintln(w, "executing /c")
	default:
		http.Error(w, "404 not found ", 404)
	}
}

type logger struct {
	inner http.Handler
}

func (l *logger) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("start %s\n", time.Now().String())
	l.inner.ServeHTTP(w, r)
	log.Printf("finish %s\n", time.Now().String())
}

// func hello(w http.ResponseWriter, r *http.Request) {
// 	fmt.Fprintf(w, "Hello %s\n", r.URL.Query().Get("name"))
// }

type trival struct {
}

func (t *trival) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	log.Println("executing trival middleware")
	next(w, r)
}

type badAuth struct {
	username string
	password string
}

func (b *badAuth) ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	username := r.URL.Query().Get("username")
	password := r.URL.Query().Get("password")
	if username != b.username || password != b.password {
		http.Error(w, "unauthorized", 401)
		return
	}
	ctx := context.WithValue(r.Context(), "username", username)
	r = r.WithContext(ctx)
	next(w, r)
}

func hello(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)
	fmt.Fprintf(w, "hi %s\n", username)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/hello", hello).Methods("GET")
	n := negroni.Classic()
	n.Use(&badAuth{username: "happyhzr", password: "password"})
	n.UseHandler(r)
	log.Fatalln(http.ListenAndServe(":8000", n))
}
