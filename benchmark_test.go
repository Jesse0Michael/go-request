package request

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
)

type BenchReq struct {
	Request struct {
		State string `json:"state"`
	} `body:"application/json"`
	User    string   `path:"user"`
	Active  bool     `query:"active"`
	Friends []string `query:"friend,explode"`
	Delay   int      `header:"X-DELAY"`
}

var expected = BenchReq{
	Request: struct {
		State string `json:"state"`
	}{State: "active"},
	User:    "adam",
	Active:  true,
	Friends: []string{"bob", "steve"},
	Delay:   60,
}

func BenchmarkDecode(b *testing.B) {
	r := benchReq()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		var req BenchReq
		err := Decode(r, &req)
		if err != nil {
			b.Fatal("failed to decode", err.Error())
		}
		if !reflect.DeepEqual(req, expected) {
			b.Errorf("Decode(r, &req) = %v, want %v", req, expected)
		}
		// reset body
		r.Body = io.NopCloser(bytes.NewReader([]byte(`{"state":"active"}`)))
	}
}

func BenchmarkBaseline(b *testing.B) {
	r := benchReq()
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		req, err := baselineDecode(r)
		if err != nil {
			b.Fatal("failed to decode", err.Error())
		}
		if !reflect.DeepEqual(*req, expected) {
			b.Errorf("baselineDecode(r) = %v, want %v", *req, expected)
		}
		// reset body
		r.Body = io.NopCloser(bytes.NewReader([]byte(`{"state":"active"}`)))
	}
}

func benchReq() *http.Request {
	url := "/users/adam?active=true&friend=bob&friend=steve"
	body := bytes.NewReader([]byte(`{"state":"active"}`))
	r := httptest.NewRequest(http.MethodPut, url, body).WithContext(context.TODO())
	r.Header.Set("X-Delay", "60")
	r.Header.Set("Content-Type", "application/json")
	vars := map[string]string{
		"user": "adam",
	}
	r = mux.SetURLVars(r, vars)
	return r
}

func baselineDecode(r *http.Request) (*BenchReq, error) {
	query := r.URL.Query()
	vars := mux.Vars(r)

	active, err := strconv.ParseBool(query.Get("active"))
	if err != nil {
		return nil, err
	}
	delay, err := strconv.Atoi(r.Header.Get("X-DELAY"))
	if err != nil {
		return nil, err
	}

	friends, _ := query["friend"]

	b, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	var data struct {
		State string `json:"state"`
	}
	err = json.Unmarshal(b, &data)
	if err != nil {
		return nil, err
	}

	return &BenchReq{
		Request: data,
		User:    vars["user"],
		Active:  active,
		Delay:   delay,
		Friends: friends,
	}, nil
}
