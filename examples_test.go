package request

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
)

func ExampleDecode() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Active bool   `query:"active"`
			State  string `json:"state"`
			Delay  int    `header:"X-DELAY"`
		}
		err := Decode(r, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err.Error())
		}

		fmt.Printf("%+v\n", req)
	}))
	defer ts.Close()

	body := `{"state":"idle"}`
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"?active=true", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DELAY", "60")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("request failed")
	}
	if resp.StatusCode == http.StatusBadRequest {
		fmt.Println("decode failed")
	}
	// Output:
	// {Active:true State:idle Delay:60}
}

func ExampleDecode_mux() {
	r := mux.NewRouter()
	r.Handle("/users/{user}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			User   string `path:"user"`
			Active bool   `query:"active"`
			State  string `json:"state"`
			Delay  int    `header:"X-DELAY"`
		}
		err := Decode(r, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err.Error())
		}

		fmt.Printf("%+v\n", req)
	}))

	body := `{"state":"idle"}`
	req, _ := http.NewRequest(http.MethodPost, "http://www.example.com/users/adam?active=true", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DELAY", "60")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code == http.StatusBadRequest {
		fmt.Println("decode failed")
	}
	// Output:
	// {User:adam Active:true State:idle Delay:60}
}

func ExampleDecode_body() {
	r := mux.NewRouter()
	r.Handle("/users/{user}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Request struct {
				State string `json:"state"`
			} `body:"application/json"`
			Active bool `query:"active"`
			Delay  int  `header:"X-DELAY"`
		}
		err := Decode(r, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Println(err.Error())
		}

		fmt.Printf("%+v\n", req)
	}))

	body := `{"state":"idle"}`
	req, _ := http.NewRequest(http.MethodPost, "http://www.example.com/users/adam?active=true", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DELAY", "60")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code == http.StatusBadRequest {
		fmt.Println("decode failed")
	}
	// Output:
	// {Request:{State:idle} Active:true Delay:60}
}
