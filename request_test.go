package request

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
)

func Test_decodeBody(t *testing.T) {
	tests := []struct {
		name    string
		r       *http.Request
		data    interface{}
		want    interface{}
		wantErr bool
	}{
		{
			name: "missing content type",
			r:    httptest.NewRequest(http.MethodPost, "/", nil),
			data: &struct{ Val string }{},
			want: &struct{ Val string }{},
		},
		{
			name: "decode json empty",
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", nil)
				r.Header.Set("Content-Type", "application/json")
				return r
			}(),
			data: &struct{ Val string }{},
			want: &struct{ Val string }{},
		},
		{
			name: "decode json",
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"Val":"success"}`))
				r.Header.Set("Content-Type", "application/json")
				return r
			}(),
			data: &struct{ Val string }{},
			want: &struct{ Val string }{Val: "success"},
		},
		{
			name: "decode json failure",
			r: func() *http.Request {
				r := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`}{`))
				r.Header.Set("Content-Type", "application/json")
				return r
			}(),
			data:    &struct{ Val string }{},
			want:    &struct{ Val string }{},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]
		t.Run(tt.name, func(t *testing.T) {
			if err := decodeBody(tt.r, tt.data); (err != nil) != tt.wantErr {
				t.Errorf("decodeBody() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(tt.data, tt.want) {
				t.Errorf("decodeBody() = %v, want %v", tt.data, tt.want)
			}
		})
	}
}
