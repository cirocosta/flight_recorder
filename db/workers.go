package db

import (
	"github.com/pkg/errors"
)

// flight_recorder_workers{state="stalled"} 3
//
func (d *Db) WorkersByState() (res map[string]float64, err error) {
	rows, err := d.db.Query(`
		SELECT
			state, count(*)
		FROM
			workers
		GROUP BY
			state
	`)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to retrieve workers by state")
		return
	}

	var (
		state string
		count float64
	)

	res = map[string]float64{}
	for rows.Next() {
		err = rows.Scan(&state, &count)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to interpret workers by state row")
			return
		}

		res[state] = count
	}

	return
}
