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
	categoryRouter := r.PathPrefix("/category").Subrouter()
	categoryRouter.Path("").Methods("GET").HandlerFunc(CategoryListHandler)
	categoryRouter.Path("").Methods("PUT").HandlerFunc(CategoryCreateHandler)
	// categoryRouter.Path("/{id}").Methods("GET").HandlerFunc(CategoryGetHandler)
	// categoryRouter.Path("/{id}").Methods("PUT").HandlerFunc(CategoryUpdateHandler)
	// categoryRouter.Path("/{id}").Methods("DELETE").HandlerFunc(CategoryDeleteHandler)
	categoryRouter.Path("/{id}/mails").Methods("GET").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.URL.Path = "/mail"
		r.URL.RawQuery = "?filter=category:" + mux.Vars(r)["id"]
		http.Redirect(w, r, r.URL.String(), http.StatusTemporaryRedirect)
	})

	log.Printf("Starting webserver on %s...", options.Listen)
	if err := http.ListenAndServe(options.Listen, nil); err != nil {
		log.Fatalf("Could not start webserver: %s", err)
	}
}
