package httpjson

import (
	"fmt"
	"net/http"
	"reflect"
)

// GetURLQuery converts a struct into a url query. The struct must have a tagname of 'query_key'
func GetURLQuery(queryParams interface{}) string {
	val := reflect.ValueOf(queryParams).Elem()
	req, _ := http.NewRequest("GET", "http://localhost", nil)
	query := req.URL.Query()

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		tag := val.Type().Field(i).Tag

		key := tag.Get("query_key")
		value := valueField.Interface()

		// If value is nil or equal to 0 do not include the query param
		if value == nil || fmt.Sprintf("%v", value) == "0" || key == "" || fmt.Sprintf("%v", value) == "" {
			continue
		}

		query.Add(key, fmt.Sprintf("%v", value))
	}

	// possible for query to be completely empty
	if query.Encode() == "" {
		return ""
	}

	// append '?' to prefix
	return "?" + query.Encode()
}
