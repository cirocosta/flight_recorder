package db

import (
	"github.com/pkg/errors"
)

func (d *Db) Teams() (datapoints []Datapoint, err error) {
	datapoints, err = d.query(`
		SELECT
			count(*)
		FROM
			teams
 	`, 0)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to count teams")
		return
	}

	return
}
