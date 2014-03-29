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

	dbMail := Mail{}
	if err := c.Find(bson.M{}).One(&dbMail); err != nil {
		t.Fatalf("Could not get new mail from DB: %s", err)
	}

	if time.Now().Sub(dbMail.Ctime) > 500*time.Millisecond {
		t.Fatalf("CTime is to far in the past: %s", dbMail.Ctime)
	}
	if time.Now().Sub(dbMail.Mtime) > 500*time.Millisecond {
		t.Fatalf("MTime is to far in the past: %s", dbMail.Ctime)
	}
	if dbMail.Status != "active" {
		t.Fatalf("Unexpected status: %s", dbMail.Status)
	}
	if dbMail.Name != "Name 0" {
		t.Fatalf("Unexpected name: %s", dbMail.Name)
	}
	if dbMail.Subject != "Some Subject" {
		t.Fatalf("Unexpected name: %s", dbMail.Subject)
	}
	if dbMail.Body != "Some Body" {
		t.Fatalf("Unexpected name: %s", dbMail.Body)
	}
	if !reflect.DeepEqual(dbMail.Category, []string{"tag1", "tag2"}) {
		t.Fatalf("Unexpected categories %#v", dbMail.Category)
	}

	newMail := Mail{}
	if err := jsonRemarshal(&newMail, resp); err != nil {
		t.Fatalf("Could not analyze response: %s", err)
	}
	if !reflect.DeepEqual(newMail, dbMail) {
		t.Fatalf("Unexpected response:\n%#v\n%#v", newMail, dbMail)
	}

	expectedURL := "http://host/mail/" + dbMail.ID.Hex()
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

func TestMailUpdateHandler(t *testing.T) {
	c := db.C(MAIL_COLLECTION)
	defer c.DropCollection()
	dataSet := insertMailTestData(c)

	r := mux.NewRouter()
	r.HandleFunc("/mail/{id}", MailUpdateHandler)

	rr := httptest.NewRecorder()
	req := mustRequest("PUT", "http://host/mail/"+dataSet[0].ID.Hex(), `
		{
			"name":"NewName",
			"category": ["tag0"],
			"subject": "NewSubject",
			"body": "NewBody"
		}`)
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}
	expectedData := dataSet[0]
	expectedData.Name = "NewName"
	expectedData.Category = []string{"tag0"}
	expectedData.Subject = "NewSubject"
	expectedData.Body = "NewBody"

	data := Mail{}
	if err := c.FindId(expectedData.ID).One(&data); err != nil {
		t.Fatalf("Could not get data: %s", err)
	}
	// Ignore mtime
	expectedData.Mtime = data.Mtime

	if !reflect.DeepEqual(data, expectedData) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func TestMailDeleteHandler(t *testing.T) {
	c := db.C(MAIL_COLLECTION)
	defer c.DropCollection()

	dataSet := insertMailTestData(c)

	r := mux.NewRouter()
	r.HandleFunc("/mail/{id}", MailDeleteHandler)

	rr := httptest.NewRecorder()
	req := mustRequest("DELETE", "http://host/mail/"+dataSet[0].ID.Hex(), "")
	r.ServeHTTP(rr, req)
	if rr.Code != http.StatusNoContent {
		t.Fatalf("Unexpected status code %d", rr.Code)
	}

	expectedData := dataSet[0]
	data := Mail{}
	if err := c.FindId(expectedData.ID).One(&data); err != nil {
		t.Fatalf("Could not get data: %s", err)
	}
	// Ignore mtime
	expectedData.Mtime = data.Mtime
	expectedData.Status = "deleted"

	if !reflect.DeepEqual(data, expectedData) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

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
