package request

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/gorilla/mux"
)

type MyRequest struct {
	Name  string `path:"name"`
	Game  string `query:"game"`
	State string `json:"state"`
	Delay string `header:"X-DELAY"`
}

func ExampleDecode() {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MyRequest
		err := Decode(r, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		fmt.Printf("%+v\n", req)
	}))
	defer ts.Close()

	body := `{"state":"idle"}`
	req, _ := http.NewRequest(http.MethodPost, ts.URL+"?game=go", strings.NewReader(body))
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
	// {Name: Game:go State:idle Delay:60}
}

func ExampleDecode_mux() {
	r := mux.NewRouter()
	r.Handle("/test/{name}", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req MyRequest
		err := Decode(r, &req)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}

		fmt.Printf("%+v\n", req)
	}))

	body := `{"state":"idle"}`
	req, _ := http.NewRequest(http.MethodPost, "http://www.example.com/test/user?game=go", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-DELAY", "60")
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	if rec.Code == http.StatusBadRequest {
		fmt.Println("decode failed")
	}
	// Output:
	// {Name:user Game:go State:idle Delay:60}
}
