package tsmapi

import "time"

// ValueStatus represents the status of a given point
// which could be valid or not valid
type ValueStatus int

const (
	VALID ValueStatus = iota
	NOT_VALID
)

type Point struct {
	Time   time.Time   `json:"time"`
	Value  float64     `json:"value"`
	Status ValueStatus `json:"state"`
}

type Points []Point
type ByTime Points

func (ps ByTime) Len() int {
	return len(ps)
}

func (ps ByTime) Swap(i, j int) {
	ps[i], ps[j] = ps[j], ps[i]
}

func (ps ByTime) Less(i, j int) bool {
	return ps[i].Time.Before(ps[j].Time)
}

// Timeseries represents a time series with a unique integer Id
// a name a list of Points and the time of the last write on this
// series.
type Timeseries struct {
	Id          int       `json:"id"`
	Name        string    `json:"name"`
	Values      Points    `json:"values"`
	LastChanged time.Time `json:"lastChanged"`
}
