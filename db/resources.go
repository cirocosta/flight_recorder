package db

import (
	"github.com/pkg/errors"
)

// flight_recorder_resources{pipeline="",team=""}
//
func (d *Db) Resources() (datapoints []Datapoint, err error) {
	datapoints, err = d.query(`
		SELECT 
			count(*), pipelines.name AS pipeline, teams.name AS team
		FROM 
			resources 
		INNER JOIN 
			pipelines 
			ON resources.pipeline_id = pipelines.id 
		INNER JOIN 
			teams 
			ON pipelines.team_id = teams.id
		GROUP BY
			pipelines.name, teams.name;
	`, 2)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to pipelines by team")
		return
	}

	return
}
