package api

import (
	"compress/gzip"
	"encoding/json"
	"net/http"
	"strings"
)

func writeJson(w http.ResponseWriter, zip bool, data interface{}) error {
	if !zip {
		return json.NewEncoder(w).Encode(data)
	}
	gw := gzip.NewWriter(w)
	return json.NewEncoder(gw).Encode(data)
}

// Timeseries represtation for testing
type Timeseries struct {
	Name string
	Id   int
}

func ReadAllTimeseries(w http.ResponseWriter, r *http.Request) {
	timeseries := Timeseries{"asdf", 1}
	zip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
	writeJson(w, zip, timeseries)
}

func GetTimeseriesInfo(w http.ResponseWriter, r *http.Request) {

}

func ReadTimeseries(w http.ResponseWriter, r *http.Request) {

}

func WriteTimeseries(w http.ResponseWriter, r *http.Request) {

}
