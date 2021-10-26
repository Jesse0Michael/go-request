package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"

	"github.com/gorilla/mux"
)

// Decode an HTTP request into the provided struct
func Decode(r *http.Request, data interface{}) error {
	typ := reflect.TypeOf(data)
	if typ == nil {
		return fmt.Errorf("invalid decode type: nil")
	}
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}
	if typ.Kind() != reflect.Struct {
		return fmt.Errorf("invalid decode type: %v", typ.Kind())
	}

	return decodeRequest(r, typ, data)
}

func decodeRequest(r *http.Request, t reflect.Type, data interface{}) error {
	query := r.URL.Query()
	vars := mux.Vars(r)
	body := false
	for i := 0; i < t.NumField(); i++ {
		typ := t.Field(i)
		val := reflect.ValueOf(data).Elem().Field(i)

		queryTag := typ.Tag.Get("query")
		if queryTag != "" {
			if query.Has(queryTag) {
				v, err := resolve(val.Interface(), query.Get(queryTag))
				if err != nil {
					return err
				}
				val.Set(reflect.ValueOf(v))
			}
		}

		pathTag := typ.Tag.Get("path")
		if pathTag != "" {
			if path, ok := vars[pathTag]; ok {
				v, err := resolve(val.Interface(), path)
				if err != nil {
					return err
				}
				val.Set(reflect.ValueOf(v))
			}
		}

		headerTag := typ.Tag.Get("header")
		if headerTag != "" {
			if r.Header.Get(headerTag) != "" {
				v, err := resolve(val.Interface(), r.Header.Get(headerTag))
				if err != nil {
					return err
				}
				val.Set(reflect.ValueOf(v))
			}
		}

		bodyTag := typ.Tag.Get("body")
		if bodyTag != "" {
			body = true
			v := reflect.New(typ.Type).Interface()
			if err := decodeBody(r, v); err != nil {
				return err
			}
			val.Set(reflect.ValueOf(v).Elem())
		}
	}
	if !body {
		err := decodeBody(r, data)
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeBody(r *http.Request, data interface{}) error {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	switch r.Header.Get("Content-Type") {
	case "application/json":
		err := json.Unmarshal(b, &data)
		if err != nil {
			return err
		}
	}

	return nil
}
