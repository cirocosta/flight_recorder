package db

import (
	"github.com/pkg/errors"
)

func (d *Db) Pipelines() (datapoints []Datapoint, err error) {
	datapoints, err = d.query(`
		SELECT
			count(*), teams.name AS team_name
		FROM
			pipelines
		INNER JOIN
			teams ON pipelines.team_id = teams.id
		GROUP BY 
			team_name
	`, 1)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to pipelines by team")
		return
	}

	return
}
