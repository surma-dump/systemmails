package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"

	"labix.org/v2/mgo/bson"
)

const (
	CATEGORY_COLLECTION = "categories"
)

type Category struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string        `bson:"name" json:"name"`
}

func CategoryListHandler(w http.ResponseWriter, r *http.Request) {
	iter, offset, count, err := FilterIter(db.C(CATEGORY_COLLECTION), r)
	if err != nil {
		http.Error(w, fmt.Sprintf("Bad Request (%s)", err), http.StatusBadRequest)
		return
	}

	result := []Category{}
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

func CategoryCreateHandler(w http.ResponseWriter, r *http.Request) {
	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	payload := Category{}
	if err := dec.Decode(&payload); err != nil {
		http.Error(w, "Bad Request (invalid payload)", http.StatusBadRequest)
		return
	}

	payload.ID = bson.NewObjectId()

	if err := db.C(CATEGORY_COLLECTION).Insert(payload); err != nil {
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

func CategoryGetHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	result := Category{}
	if err := db.C(CATEGORY_COLLECTION).FindId(bson.ObjectIdHex(id)).One(&result); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func CategoryUpdateHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	payload := Category{}
	if err := dec.Decode(&payload); err != nil {
		http.Error(w, "Bad Request (invalid payload)", http.StatusBadRequest)
		return
	}
	payload.ID = bson.ObjectIdHex(id)
	if err := db.C(CATEGORY_COLLECTION).UpdateId(payload.ID, payload); err != nil {
		log.Printf("Could not update value: %s", err)
		http.Error(w, "Internal server error (could not update value)", http.StatusInternalServerError)
		return
	}
	http.Error(w, "", http.StatusNoContent)
}

func CategoryDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	if err := db.C(CATEGORY_COLLECTION).RemoveId(bson.ObjectIdHex(id)); err != nil {
		log.Printf("Could not remove value: %s", err)
		http.Error(w, "Internal server error (could not remove value)", http.StatusInternalServerError)
		return
	}
	http.Error(w, "", http.StatusNoContent)
}
