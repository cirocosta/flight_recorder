package collectors

import (
	"github.com/cirocosta/flight_recorder/db"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	containersDescription = prometheus.NewDesc(
		"flight_recorder_containers",
		"Number of containers",
		[]string{"type", "worker", "state"},
		nil,
	)
)

type Containers struct {
	Db *db.Db
}

func (w *Containers) Describe(ch chan<- *prometheus.Desc) {
	ch <- containersDescription
}

func (w *Containers) Collect(ch chan<- prometheus.Metric) {
	datapoints, err := w.Db.Containers()
	if err != nil {
		panic(err)
	}

	for _, datapoint := range datapoints {
		ch <- prometheus.MustNewConstMetric(
			containersDescription,
			prometheus.UntypedValue,
			datapoint.Value,
			datapoint.LabelSet...,
		)
	}
}
