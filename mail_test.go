package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"labix.org/v2/mgo"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"labix.org/v2/mgo/bson"
)

func TestMailCreateHandler(t *testing.T) {
	c := db.C(MAIL_COLLECTION)
	defer c.DropCollection()

	rr := httptest.NewRecorder()
	MailCreateHandler(rr, mustRequest("POST", "http://host/mail", `
		{
			"name": "Name 0",
			"subject": "Some Subject",
			"body": "Some Body",
			"category": ["tag1", "tag2"]
		}`))
	if rr.Code != http.StatusCreated {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}

	resp := bson.M{}
	if err := json.Unmarshal(rr.Body.Bytes(), &resp); err != nil {
		t.Fatalf("Could not unmarshal response: %s", err)
	}

	newMail := Mail{}
	if err := c.Find(bson.M{}).One(&newMail); err != nil {
		t.Fatalf("Could not get new mail from DB: %s", err)
	}

	if time.Now().Sub(newMail.Ctime) > 500*time.Millisecond {
		t.Fatalf("CTime is to far in the past: %s", newMail.Ctime)
	}
	if time.Now().Sub(newMail.Mtime) > 500*time.Millisecond {
		t.Fatalf("MTime is to far in the past: %s", newMail.Ctime)
	}
	if newMail.Status != "active" {
		t.Fatalf("Unexpected status: %s", newMail.Status)
	}
	if newMail.Name != "Name 0" {
		t.Fatalf("Unexpected name: %s", newMail.Name)
	}
	if newMail.Subject != "Some Subject" {
		t.Fatalf("Unexpected name: %s", newMail.Subject)
	}
	if newMail.Body != "Some Body" {
		t.Fatalf("Unexpected name: %s", newMail.Body)
	}
	if !reflect.DeepEqual(newMail.Category, []string{"tag1", "tag2"}) {
		t.Fatalf("Unexpected categories %#v", newMail.Category)
	}

	expectedURL := "http://host/mail/" + newMail.ID.Hex()
	if loc := rr.Header().Get("Location"); loc != expectedURL {
		t.Fatalf("Unexpected Location header %s", loc)
	}
	if resp["url"] != expectedURL {
		t.Fatalf("Unexpected URl values: %s", resp["url"])
	}
}

func TestMailGetHandler(t *testing.T) {
	c := db.C(MAIL_COLLECTION)
	defer c.DropCollection()
	dataSet := insertMailTestData(c)

	r := mux.NewRouter()
	r.HandleFunc("/mail/{id}", MailGetHandler)

	rr := httptest.NewRecorder()
	req := mustRequest("GET", "http://host/mail/"+dataSet[0].ID.Hex(), "")
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}

	data := Mail{}
	if err := json.Unmarshal(rr.Body.Bytes(), &data); err != nil {
		t.Fatalf("Could not unmarshal response: %s", err)
	}

	if !reflect.DeepEqual(data, dataSet[0]) {
		t.Fatalf("Unexpected data:\n%#v\n%#v", data, dataSet[0])
	}
}

/*
func TestMailUpdateHandler(t *testing.T) {
	c := db.C(CATEGORY_COLLECTION)
	defer c.DropCollection()

	dataSet := []Category{
		{
			ID:   bson.NewObjectId(),
			Name: "Name0",
		},
		{
			ID:   bson.NewObjectId(),
			Name: "Name1",
		},
	}

	for _, d := range dataSet {
		if err := c.Insert(d); err != nil {
			t.Fatalf("Could not insert test data: %s", err)
		}
	}

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

func TestMailDeleteHandler(t *testing.T) {
	c := db.C(CATEGORY_COLLECTION)
	defer c.DropCollection()

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
			t.Fatalf("Could not insert test data: %s", err)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/category/{id}", CategoryDeleteHandler)

	rr := httptest.NewRecorder()
	req := mustRequest("DELETE", "http://host/category/"+dataSet[0].ID.Hex(), "")
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}
	dataSet[0].Name = "NewName"

	data := []Category{}
	if err := c.Find(bson.M{}).All(&data); err != nil {
		t.Fatalf("Could not get data: %s", err)
	}

	if !reflect.DeepEqual(data, dataSet[1:]) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}
*/

func insertMailTestData(c *mgo.Collection) []Mail {
	// MongoDB doesnt store Nanoseconds
	t := time.Now().Round(time.Millisecond)
	dataSet := []Mail{
		{
			ID:       bson.NewObjectId(),
			Name:     "Name 0",
			Category: []string{"tag1", "tag2"},
			Author:   "Icke",
			Ctime:    t,
			Mtime:    t,
			Status:   "active",
			Subject:  "Subject",
			Body:     "Body",
		},
		{
			ID:       bson.NewObjectId(),
			Name:     "Name 1",
			Category: []string{"tag2", "tag3"},
			Author:   "Icke",
			Ctime:    t,
			Mtime:    t,
			Status:   "active",
			Subject:  "Subject",
			Body:     "Body",
		},
		{
			ID:       bson.NewObjectId(),
			Name:     "Name 0",
			Category: []string{"tag3", "tag4"},
			Author:   "Icke",
			Ctime:    t,
			Mtime:    t,
			Status:   "deleted",
			Subject:  "Subject",
			Body:     "Body",
		},
	}

	for _, d := range dataSet {
		if err := c.Insert(d); err != nil {
			panic(fmt.Sprintf("Could not insert test data: %s", err))
		}
	}
	return dataSet
}
