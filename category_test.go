package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func init() {
	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	db = session.DB("")
}

type categoryResult struct {
	Offset int        `json:"offset"`
	Count  int        `json:"count"`
	Result []Category `json:"result"`
}

func TestCategoryCreateHandler(t *testing.T) {
	c := db.C(CATEGORY_COLLECTION)
	defer c.DropCollection()

	rr := httptest.NewRecorder()
	CategoryCreateHandler(rr, mustRequest("POST", "http://host/category", `{"name": "Name 0"}`))
	if rr.Code != http.StatusCreated {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}

	data := bson.M{}
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("Could not unmarshal response: %s", err)
	}

	newCategory := bson.M{}
	if err := c.Find(bson.M{}).One(&newCategory); err != nil {
		t.Fatalf("Could not get new category from DB: %s", err)
	}

	expectedData := bson.M{
		"_id":  newCategory["_id"].(bson.ObjectId).Hex(),
		"name": "Name 0",
		"url":  "http://host/category/" + newCategory["_id"].(bson.ObjectId).Hex(),
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Fatalf("Unexpected data %#v", data)
	}
	if loc := rr.Header().Get("Location"); loc != expectedData["url"] {
		t.Fatalf("Unexpected Location header %s", loc)
	}
}

func TestCategoryGetHandler(t *testing.T) {
	c := db.C(CATEGORY_COLLECTION)
	defer c.DropCollection()

	dataSet := insertCategoryTestData(c)

	r := mux.NewRouter()
	r.HandleFunc("/category/{id}", CategoryGetHandler)

	rr := httptest.NewRecorder()
	req := mustRequest("GET", "http://host/category/"+dataSet[0].ID.Hex(), "")
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}

	data := Category{}
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("Could not unmarshal response: %s", err)
	}

	if !reflect.DeepEqual(data, dataSet[0]) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func TestCategoryUpdateHandler(t *testing.T) {
	c := db.C(CATEGORY_COLLECTION)
	defer c.DropCollection()

	dataSet := insertCategoryTestData(c)

	r := mux.NewRouter()
	r.HandleFunc("/category/{id}", CategoryUpdateHandler)

	rr := httptest.NewRecorder()
	req := mustRequest("PUT", "http://host/category/"+dataSet[0].ID.Hex(), `{"name":"NewName"}`)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}
	dataSet[0].Name = "NewName"

	data := []Category{}
	if err := c.Find(bson.M{}).All(&data); err != nil {
		t.Fatalf("Could not get data: %s", err)
	}

	if !reflect.DeepEqual(data, dataSet) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func TestCategoryDeleteHandler(t *testing.T) {
	c := db.C(CATEGORY_COLLECTION)
	defer c.DropCollection()

	dataSet := insertCategoryTestData(c)

	r := mux.NewRouter()
	r.HandleFunc("/category/{id}", CategoryDeleteHandler)

	rr := httptest.NewRecorder()
	req := mustRequest("DELETE", "http://host/category/"+dataSet[0].ID.Hex(), "")
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}

	data := Category{}
	if err := c.FindId(dataSet[0].ID).One(&data); err != mgo.ErrNotFound {
		t.Fatalf("Dataset still present: %s", err)
	}
}

func insertCategoryTestData(c *mgo.Collection) []Category {
	dataSet := []Category{
		{
			ID:   bson.NewObjectId(),
			Name: "Name0",
		},
		{
			ID:   bson.NewObjectId(),
			Name: "Name1",
		},
		{
			ID:   bson.NewObjectId(),
			Name: "Name2",
		},
	}

	for _, d := range dataSet {
		if err := c.Insert(d); err != nil {
			panic(fmt.Sprintf("Could not insert test data: %s", err))
		}
	}
	return dataSet
}
