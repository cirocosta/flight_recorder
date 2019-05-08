package collectors

import (
	"github.com/cirocosta/flight_recorder/db"
	"github.com/prometheus/client_golang/prometheus"
)

type Collector struct {
	Description   *prometheus.Desc
	RetrievalFunc func() (datapoints []db.Datapoint, err error)
}

func (c *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.Description
}

func (c *Collector) Collect(ch chan<- prometheus.Metric) {
	datapoints, err := c.RetrievalFunc()
	if err != nil {
		panic(err)
	}

	for _, datapoint := range datapoints {
		ch <- prometheus.MustNewConstMetric(
			c.Description,
			prometheus.UntypedValue,
			datapoint.Value,
			datapoint.LabelSet...,
		)
	}
}
