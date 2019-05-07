package collectors

import (
	"github.com/cirocosta/flight_recorder/db"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	workersByStateDescription = prometheus.NewDesc(
		"flight_recorder_workers",
		"Per-state worker count",
		[]string{"state"},
		nil,
	)
)

type WorkersByState struct {
	Db *db.Db
}

func (w *WorkersByState) Describe(ch chan<- *prometheus.Desc) {
	ch <- workersByStateDescription
}

func (w *WorkersByState) Collect(ch chan<- prometheus.Metric) {
	workersByState, err := w.Db.WorkersByState()
	if err != nil {
		panic(err)
	}

	for state, count := range workersByState {
		ch <- prometheus.MustNewConstMetric(
			workersByStateDescription,
			prometheus.UntypedValue,
			count,
			state,
		)
	}
}