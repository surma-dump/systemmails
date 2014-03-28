package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

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
	query := bson.M{}
	container := bson.M{}
	if v := r.FormValue("filter"); v != "" {
		fields := strings.Split(v, ":")
		if len(fields) != 2 {
			http.Error(w, "Bad Request (invalid filter)", http.StatusBadRequest)
			return
		}

		query[fields[0]] = fields[1]
	}
	qry := db.C(CATEGORY_COLLECTION).Find(query)
	if v := r.FormValue("skip"); v != "" {
		skip, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			http.Error(w, "Bad Request (invalid skip value)", http.StatusBadRequest)
			return
		}
		container["offset"] = skip
		qry.Skip(int(skip))
	}
	if v := r.FormValue("limit"); v != "" {
		limit, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			http.Error(w, "Bad Request (invalid limit value)", http.StatusBadRequest)
			return
		}
		qry.Limit(int(limit))
	}

	container["count"], _ = qry.Count()
	result := []Category{}
	if err := qry.All(&result); err != nil {
		log.Printf("Could not fetch values: %s", err)
		http.Error(w, "Internal server error (could not fetch values)", http.StatusInternalServerError)
		return
	}
	container["result"] = result

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(container)
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

/*
func (mcs *MgoCollectionServer) GetHandler(w http.ResponseWriter, r *http.Request) {
	pathElems := strings.Split(r.URL.Path, "/")
	id := pathElems[len(pathElems)-1]

	var result interface{}
	if err := mcs.Collection.FindId(id).One(&result); err != nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	enc := json.NewEncoder(w)
	enc.Encode(result)
}

func (mcs *MgoCollectionServer) PutHandler(w http.ResponseWriter, r *http.Request) {
	pathElems := strings.Split(r.URL.Path, "/")
	id := pathElems[len(pathElems)-1]

	dec := json.NewDecoder(r.Body)
	defer r.Body.Close()

	var payload interface{}
	if err := dec.Decode(&payload); err != nil {
		http.Error(w, "Bad Request (invalid payload)", http.StatusBadRequest)
		return
	}
	if mcs.Validate != nil {
		if err := mcs.Validate(payload); err != nil {
			http.Error(w, "Bad Request (validation failed)", http.StatusBadRequest)
			return
		}
	}

	if err := mcs.Collection.UpdateId(id, payload); err != nil {
		log.Printf("Could not update value: %s", err)
		http.Error(w, "Internal server error (could not update value)", http.StatusInternalServerError)
		return
	}
	http.Error(w, "", http.StatusNoContent)
}
*/
