package main

import (
	"encoding/json"
	"fmt"
	// "github.com/gorilla/mux"
	"log"
	"net/http"
	"time"

	"labix.org/v2/mgo/bson"
)

const (
	MAIL_COLLECTION = "mails"
)

type Mail struct {
	ID       bson.ObjectId `bson:"_id,omitempty"    json:"_id,omitempty"`
	Name     string        `bson:"name"             json:"name"`
	Category []string      `bson:"category"         json:"category"`
	Author   string        `bson:"author,omitempty" json:"author,omitempty"`
	Ctime    time.Time     `bson:"ctime,omitempty"  json:"ctime,omitempty"`
	Mtime    time.Time     `bson:"mtime,omitempty"  json:"mtime,omitempty"`
	Status   string        `bson:"status,omitempty" json:"status,omitempty"`
	Subject  interface{}   `bson:"subject"          json:"subject"`
	Body     interface{}   `bson:"body"             json:"body"`
}

func MailListHandler(w http.ResponseWriter, r *http.Request) {
	iter, offset, count, err := FilterIter(db.C(MAIL_COLLECTION), r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request (%s)", err), http.StatusBadRequest)
		return
	}

	result := []Mail{}
	if err := iter.All(&result); err != nil {
		log.Printf("Could not fetch values: %s", err)
		http.Error(w, "Internal server error (could not fetch values)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bson.M{
		"count":  count,
		"offset": offset,
		"result": result,
	})
}

/*
func MailCreateHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	payload := Category{}
	if err := dec.Decode(&payload); err != nil {
		http.Error(w, "Bad Request (invalid payload)", http.StatusBadRequest)
		return
	}

	payload.ID = bson.NewObjectId()

	if err := db.C(MAIL_COLLECTION).Insert(payload); err != nil {
		log.Printf("Could not insert value: %s", err)
		http.Error(w, "Internal server error (could not insert value)", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	r.URL.Path += "/" + payload.ID.Hex()
	json.NewEncoder(w).Encode(bson.M{
		"_id":  payload.ID,
		"name": payload.Name,
		"url":  r.URL.String(),
	})
}

func MailGetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result := Category{}
	if err := db.C(MAIL_COLLECTION).FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func MailUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	payload := Category{}
	if err := dec.Decode(&payload); err != nil {
		http.Error(w, "Bad Request (invalid payload)", http.StatusBadRequest)
		return
	}
	payload.ID = bson.ObjectIdHex(id)
	if err := db.C(MAIL_COLLECTION).UpdateId(payload.ID, payload); err != nil {
		log.Printf("Could not update value: %s", err)
		http.Error(w, "Internal server error (could not update value)", http.StatusInternalServerError)
		return
	}
	http.Error(w, "", http.StatusNoContent)
}

func MailDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := db.C(MAIL_COLLECTION).RemoveId(bson.ObjectIdHex(id)); err != nil {
		log.Printf("Could not remove value: %s", err)
		http.Error(w, "Internal server error (could not remove value)", http.StatusInternalServerError)
		return
	}
	http.Error(w, "", http.StatusNoContent)
}
*/
