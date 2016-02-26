package tsmapi

import "time"

//
type TimeseriesStore interface {
	Read(id int, start, end time.Time) (Timeseries, error)
	Write(id int, values Points) error
}
