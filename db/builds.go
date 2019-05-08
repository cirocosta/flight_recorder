package db

import (
	"github.com/pkg/errors"
)

func (d *Db) Builds() (datapoints []Datapoint, err error) {
	datapoints, err = d.query(`
		SELECT
			count(*), status, teams.name AS team, pipelines.name AS pipeline
		FROM
			builds
		INNER JOIN
			teams
			ON teams.id = builds.team_id
		INNER JOIN
			pipelines
			ON pipelines.id = builds.pipeline_id
		GROUP BY 
			status, team, pipeline;
	`, 3)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to pipelines by team")
		return
	}

	return
}
