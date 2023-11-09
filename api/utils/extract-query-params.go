package utils

import (
	"encoding/json"
	"net/http"
)

func ExtracteQueryParams(req *http.Request) ([]byte, error) {
	queryParams := req.URL.Query()
	query := make(map[string]interface{})
	for name, values := range queryParams {
		query[name] = values[0]
	}
	return json.Marshal(query)
}
