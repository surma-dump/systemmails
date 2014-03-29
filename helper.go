package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
)

func FilterIter(c *mgo.Collection, r *http.Request) (iter *mgo.Iter, offset, count int, err error) {
	query := bson.M{}
	if v := r.FormValue("filter"); v != "" {
		fields := strings.Split(v, ":")
		if len(fields) != 2 {
			return nil, 0, 0, fmt.Errorf("Invalid filter")
		}
		query[fields[0]] = fields[1]
	}
	qry := c.Find(query)
	count, err = qry.Count()
	if err != nil {
		return nil, 0, 0, err
	}
	if v := r.FormValue("skip"); v != "" {
		skip, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("Invalid skip value")
		}
		offset = int(skip)
		qry.Skip(offset)
	}
	if v := r.FormValue("limit"); v != "" {
		limit, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("Invalid limit value")
		}
		qry.Limit(int(limit))
	}
	return qry.Iter(), offset, count, nil
}

func mustRequest(method, urlStr, body string) *http.Request {
	r, err := http.NewRequest(method, urlStr, strings.NewReader(body))
	if err != nil {
		panic(err)
	}
	return r
}

func jsonRemarshal(new interface{}, old interface{}) error {
	data, err := json.Marshal(old)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, new)
}
