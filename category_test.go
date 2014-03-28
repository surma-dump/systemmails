package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
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

func TestCategory_ListHandler_Listing(t *testing.T) {
	c := db.C(CATEGORY_COLLECTION)
	defer c.DropCollection()

	dataSet := []Category{
		{
			ID:   bson.NewObjectId(),
			Name: "Name 1",
		},
		{
			ID:   bson.NewObjectId(),
			Name: "Name 2",
		},
	}

	for _, d := range dataSet {
		if err := c.Insert(d); err != nil {
			t.Fatalf("Could not insert test data: %s", err)
		}
	}

	rr := httptest.NewRecorder()
	CategoryListHandler(rr, mustRequest("GET", "http://host/category", ""))

	if rr.Code != http.StatusOK {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}

	data := categoryResult{}
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("Could not unmarshal response: %s", err)
	}

	expectedData := categoryResult{
		Offset: 0,
		Count:  2,
		Result: dataSet,
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func TestCategory_ListHandler_Filtering(t *testing.T) {
	c := db.C(CATEGORY_COLLECTION)
	defer c.DropCollection()

	dataSet := []Category{
		{
			ID:   bson.NewObjectId(),
			Name: "Name1",
		},
		{
			ID:   bson.NewObjectId(),
			Name: "Name2",
		},
		{
			ID:   bson.NewObjectId(),
			Name: "Name3",
		},
	}

	for _, d := range dataSet {
		if err := c.Insert(d); err != nil {
			t.Fatalf("Could not insert test data: %s", err)
		}
	}

	rr := httptest.NewRecorder()
	CategoryListHandler(rr, mustRequest("GET", "http://host/category?filter=name:"+dataSet[1].Name, ""))
	if rr.Code != http.StatusOK {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}

	data := categoryResult{}
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("Could not unmarshal response: %s", err)
	}

	expectedData := categoryResult{
		Offset: 0,
		Count:  1,
		Result: dataSet[1:2],
	}
	if !reflect.DeepEqual(data, expectedData) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func TestCategory_CreateHandler(t *testing.T) {
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
}

func mustRequest(method, urlStr, body string) *http.Request {
	r, err := http.NewRequest(method, urlStr, strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	return r
}
