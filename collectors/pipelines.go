package collectors

import (
	"github.com/cirocosta/flight_recorder/db"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	pipelinesPerTeamDescription = prometheus.NewDesc(
		"flight_recorder_pipelines",
		"Number of pipelines set",
		[]string{"team"},
		nil,
	)
)

type PipelinesPerTeam struct {
	Db *db.Db
}

func (w *PipelinesPerTeam) Describe(ch chan<- *prometheus.Desc) {
	ch <- pipelinesPerTeamDescription
}

func (w *PipelinesPerTeam) Collect(ch chan<- prometheus.Metric) {
	datapoints, err := w.Db.PipelinesPerTeam()
	if err != nil {
		panic(err)
	}

	for _, datapoint := range datapoints {
		ch <- prometheus.MustNewConstMetric(
			pipelinesPerTeamDescription,
			prometheus.UntypedValue,
			datapoint.Value,
			datapoint.LabelSet...,
		)
	}
}
