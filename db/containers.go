package db

import (
	"github.com/pkg/errors"
)

func (d *Db) Containers() (datapoints []Datapoint, err error) {
	datapoints, err = d.query(`
		SELECT 
			count(*), meta_type, worker_name, state 
		FROM 
			containers 
		GROUP BY 
			meta_type, worker_name, state;
	`, 3)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to retrieve containers")
		return
	}

	return
}
