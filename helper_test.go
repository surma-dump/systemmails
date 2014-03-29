package main

import (
	"fmt"
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
	db = session.DB("systemmails-test")
}

type test struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	Name string        `bson:"name" json:"name"`
	Tags []string      `bson:"tags" json:"tags"`
}

func TestFilterIter_Data(t *testing.T) {
	c := db.C("someCollection")
	defer c.DropCollection()
	dataSet := insertHelperTestData(c)

	iter, offset, count, err := FilterIter(c, mustRequest("GET", "/", ""))
	if err != nil {
		t.Fatalf("Query failed: %s", err)
	}

	data := []test{}
	if err := iter.All(&data); err != nil {
		t.Fatalf("Obtaining data failed: %s", err)
	}
	if offset != 0 {
		t.Fatalf("Unexpected offset %d", offset)
	}
	if count != len(dataSet) {
		t.Fatalf("Unexpected data count %d", count)
	}
	if !reflect.DeepEqual(data, dataSet) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func TestFilterIter_Limit(t *testing.T) {
	c := db.C("someCollection")
	defer c.DropCollection()
	dataSet := insertHelperTestData(c)

	iter, offset, count, err := FilterIter(c, mustRequest("GET", "/?limit=2", ""))
	if err != nil {
		t.Fatalf("Query failed: %s", err)
	}

	data := []test{}
	if err := iter.All(&data); err != nil {
		t.Fatalf("Obtaining data failed: %s", err)
	}
	if offset != 0 {
		t.Fatalf("Unexpected offset %d", offset)
	}
	if count != len(dataSet) {
		t.Fatalf("Unexpected data count %d", count)
	}
	if !reflect.DeepEqual(data, dataSet[0:2]) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func TestFilterIter_Skip(t *testing.T) {
	c := db.C("someCollection")
	defer c.DropCollection()
	dataSet := insertHelperTestData(c)

	iter, offset, count, err := FilterIter(c, mustRequest("GET", "/?skip=1", ""))
	if err != nil {
		t.Fatalf("Query failed: %s", err)
	}

	data := []test{}
	if err := iter.All(&data); err != nil {
		t.Fatalf("Obtaining data failed: %s", err)
	}
	if offset != 1 {
		t.Fatalf("Unexpected offset %d", offset)
	}
	if count != len(dataSet) {
		t.Fatalf("Unexpected data count %d", count)
	}
	if !reflect.DeepEqual(data, dataSet[1:3]) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func TestFilterIter_Filter(t *testing.T) {
	c := db.C("someCollection")
	defer c.DropCollection()
	dataSet := insertHelperTestData(c)

	iter, offset, count, err := FilterIter(c, mustRequest("GET", "/?filter=tags:tag3", ""))
	if err != nil {
		t.Fatalf("Query failed: %s", err)
	}

	data := []test{}
	if err := iter.All(&data); err != nil {
		t.Fatalf("Obtaining data failed: %s", err)
	}
	if offset != 0 {
		t.Fatalf("Unexpected offset %d", offset)
	}
	if count != 2 {
		t.Fatalf("Unexpected data count %d", count)
	}
	if !reflect.DeepEqual(data, dataSet[1:3]) {
		t.Fatalf("Unexpected data: %#v", data)
	}
}

func insertHelperTestData(c *mgo.Collection) []test {
	dataSet := []test{
		{
			ID:   bson.NewObjectId(),
			Name: "Name 1",
			Tags: []string{"tag1", "tag2"},
		},
		{
			ID:   bson.NewObjectId(),
			Name: "Name 2",
			Tags: []string{"tag2", "tag3"},
		},
		{
			ID:   bson.NewObjectId(),
			Name: "Name 3",
			Tags: []string{"tag3"},
		},
	}

	for _, d := range dataSet {
		if err := c.Insert(d); err != nil {
			panic(fmt.Sprintf("Could not insert test data: %s", err))
		}
	}
	return dataSet
}
