package db

import (
	"github.com/pkg/errors"
)

func (d *Db) Workers() (datapoints []Datapoint, err error) {
	datapoints, err = d.query(`
		SELECT
			count(*), state
		FROM
			workers
		GROUP BY
			state
	`, 1)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to query workers by state")
		return
	}

	return
}
