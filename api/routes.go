package api

import "github.com/espang/router"

var routes router.Routes

func init() {
	//create all the routes for the api
	routes = router.Routes{
		{
			"ReadAllTimeseries",
			"GET",
			"/timeseries",
			ReadAllTimeseries,
		},
		{
			"GetTimeseriesInfo",
			"GET",
			"/info/{tsid}",
			GetTimeseriesInfo,
		},
		{
			"ReadTimeseries",
			"GET",
			"/data/{from}:{to}",
			ReadTimeseries,
		},
		{
			"WriteTimeseries",
			"POST",
			"/timeseries",
			WriteTimeseries,
		},
	}
}
