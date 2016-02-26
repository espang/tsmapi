package tsmapi

import (
	"compress/gzip"
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"
)

var gzipWriterPool = sync.Pool{
	New: func() interface{} { return gzip.NewWriter(nil) },
}

func makeResponse(w http.ResponseWriter, state int, msg string) {
	w.WriteHeader(state)
	if err := json.NewEncoder(w).Encode(msg); err != nil {
		log.Printf("Error: Creating error response with state %d: %s", state, err)
	}
}

func writeJson(w http.ResponseWriter, zip bool, data interface{}) {
	buf, err := json.Marshal(data)
	if err != nil {
		makeResponse(w, http.StatusInternalServerError, err.Error())
		return
	}
	if !zip {
		w.Write(buf)
		return
	}
	gzw := gzipWriterPool.Get().(*gzip.Writer)
	gzw.Reset(w)
	w.Header().Set("Content-Encoding", "gzip")
	gzw.Write(buf)
	gzw.Close()
	gzipWriterPool.Put(gzw)
	w.Header().Set("Content-Encoding", "gzip")
}

func ReadAllTimeseries(w http.ResponseWriter, r *http.Request) {
	timeseries := Timeseries{
		Id:          1,
		Name:        "asdf",
		Values:      Points{},
		LastChanged: time.Now(),
	}
	//zip := strings.Contains(r.Header.Get("Accept-Encoding"), "gzip")
	zip := false
	writeJson(w, zip, timeseries)
}

func GetTimeseriesInfo(w http.ResponseWriter, r *http.Request) {

}

func ReadTimeseries(w http.ResponseWriter, r *http.Request) {

}

func WriteTimeseries(w http.ResponseWriter, r *http.Request) {

}
