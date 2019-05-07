package collectors

import (
	"github.com/cirocosta/flight_recorder/db"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	teamsDescription = prometheus.NewDesc(
		"flight_recorder_teams",
		"Number of teams",
		nil,
		nil,
	)
)

type Teams struct {
	Db *db.Db
}

func (w *Teams) Describe(ch chan<- *prometheus.Desc) {
	ch <- teamsDescription
}

func (w *Teams) Collect(ch chan<- prometheus.Metric) {
	datapoints, err := w.Db.Teams()
	if err != nil {
		panic(err)
	}

	for _, datapoint := range datapoints {
		ch <- prometheus.MustNewConstMetric(
			teamsDescription,
			prometheus.UntypedValue,
			datapoint.Value,
			datapoint.LabelSet...,
		)
	}
}
