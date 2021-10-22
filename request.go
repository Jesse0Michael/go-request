package request

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

func Decode(r *http.Request, data interface{}) error {
	err := decodeBody(r, data)
	if err != nil {
		return err
	}

	err = decodeRequest(r, data)
	if err != nil {
		return err
	}

	return nil
}

func decodeRequest(r *http.Request, data interface{}) error {
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

	val := reflect.ValueOf(data).Elem()
	for i := 0; i < typ.NumField(); i++ {
		err := decodeField(r, typ.Field(i), val.Field(i))
		if err != nil {
			return err
		}
	}
	return nil
}

func decodeField(r *http.Request, typ reflect.StructField, val reflect.Value) error {
	query := r.URL.Query()
	queryTag := typ.Tag.Get("query")
	if queryTag != "" {
		if query.Has(queryTag) {
			v := query.Get(queryTag)
			val.Set(reflect.ValueOf(v))
		}
	}

	headerTag := typ.Tag.Get("header")
	if headerTag != "" {
		if r.Header.Get(headerTag) != "" {
			v := r.Header.Get(headerTag)
			val.Set(reflect.ValueOf(v))
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
