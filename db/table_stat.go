package db

import (
	"strings"

	"github.com/pkg/errors"
)

var measurements = []string{
	"seq_scan",
	"seq_tup_read",
	"idx_scan",
	"idx_tup_fetch",
	"n_tup_ins",
	"n_tup_upd",
	"n_tup_del",
	"n_tup_hot_upd",
	"n_live_tup",
	"n_dead_tup",
}

func tableStatQuery() string {
	columns := "relname"

	for _, measurement := range measurements {
		columns += ","
		if strings.HasPrefix(measurement, "idx_") {
			columns += "coalesce(" + measurement + ", 0)"
		} else {
			columns += measurement
		}
	}

	return `SELECT ` + columns + ` FROM pg_stat_user_tables;`
}

// TableStat retrieves statistics from `pg_stat_user_tables`.
//
// 	- seq_scan		- seq scans initiated on this table
// 	- seq_tup_read		- live rows fetched by sequential scans
// 	- idx_scan		- index scans initiated on this table
// 	- idx_tup_fetch		- live rows fetched by index scans
// 	- n_tup_ins		- rows inserted
// 	- n_tup_upd 		- rows updated
// 	- n_tup_del 		- rows deleted
// 	- n_tup_hot_upd 	- rows HOT updated
// 	- n_live_tup 		- lives rows
// 	- n_dead_tup		- read rows
//
func (d *Db) TableStat() (datapoints []Datapoint, err error) {
	rows, err := d.db.Query(tableStatQuery())

	if err != nil {
		err = errors.Wrapf(err,
			"failed to execute `pg_stat_user_tables` query")
		return
	}

	var (
		tableName       string = ""
		valuesHolder           = make([]float64, len(measurements))
		scanningTargets        = []interface{}{&tableName}
	)

	for idx := range measurements {
		scanningTargets = append(scanningTargets, &valuesHolder[idx])
	}

	for rows.Next() {
		err = rows.Scan(scanningTargets...)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to scan pg_stat_user_tables row")
			return
		}

		for idx, measurement := range measurements {
			datapoints = append(datapoints, Datapoint{
				LabelSet: []string{
					measurement,
					tableName,
				},
				Value: valuesHolder[idx],
			})
		}
	}

	return
}
