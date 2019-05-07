package collectors

import (
	"github.com/cirocosta/flight_recorder/db"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	resourcesDescription = prometheus.NewDesc(
		"flight_recorder_resources",
		"Number of resources configured",
		[]string{"pipeline", "team"},
		nil,
	)
)

type Resources struct {
	Db *db.Db
}

func (w *Resources) Describe(ch chan<- *prometheus.Desc) {
	ch <- resourcesDescription
}

func (w *Resources) Collect(ch chan<- prometheus.Metric) {
	datapoints, err := w.Db.Resources()
	if err != nil {
		panic(err)
	}

	for _, datapoint := range datapoints {
		ch <- prometheus.MustNewConstMetric(
			resourcesDescription,
			prometheus.UntypedValue,
			datapoint.Value,
			datapoint.LabelSet...,
		)
	}
}
