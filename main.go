package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/voxelbrain/goptions"
	"labix.org/v2/mgo"
)

var (
	options = struct {
		Listen  string        `goptions:"-l, --listen, description='Address to bind to'"`
		MongoDB string        `goptions:"-m, --mongodb, description='URL of the MongoDB', obligatory"`
		Help    goptions.Help `goptions:"-h, --help, description='Show this help'"`
	}{
		Listen: "localhost:8000",
	}
)

var db *mgo.Database

func main() {
	goptions.ParseAndFail(&options)

	session, err := mgo.Dial(options.MongoDB)
	if err != nil {
		log.Fatalf("Could not connect to MongoDB: %s", err)
	}
	db = session.DB("") // Use database from URL

	r := mux.NewRouter()

	r.Path("/category").Methods("GET").HandlerFunc(CategoryListHandler)
	r.Path("/category").Methods("POST").HandlerFunc(CategoryCreateHandler)
	r.Path("/category/{id}").Methods("GET").HandlerFunc(CategoryGetHandler)
	r.Path("/category/{id}").Methods("PUT").HandlerFunc(CategoryUpdateHandler)
	r.Path("/category/{id}").Methods("DELETE").HandlerFunc(CategoryDeleteHandler)
	r.Path("/category/{id}/mails").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/mail"
		r.URL.RawQuery = "?filter=category:" + mux.Vars(r)["id"]
		http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
	})

	r.Path("/mail").Methods("GET").HandlerFunc(MailListHandler)
	r.Path("/mail").Methods("POST").HandlerFunc(MailCreateHandler)
	r.Path("/mail/{id}").Methods("GET").HandlerFunc(MailGetHandler)
	r.Path("/mail/{id}").Methods("PUT").HandlerFunc(MailUpdateHandler)
	r.Path("/mail/{id}").Methods("DELETE").HandlerFunc(CategoryDeleteHandler)

	injectHostname := func(w http.ResponseWriter, req *http.Request) {
		req.URL.Host = req.Host
		req.URL.Scheme = "http"
		r.ServeHTTP(w, req)
	}

	log.Printf("Starting webserver on %s...", options.Listen)
	if err := http.ListenAndServe(options.Listen, http.HandlerFunc(injectHostname)); err != nil {
		log.Fatalf("Could not start webserver: %s", err)
	}
}
