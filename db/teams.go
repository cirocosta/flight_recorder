package db

// import (
// 	"github.com/pkg/errors"
// )

// flight_recorder_teams_total
//
// func (d *Db) Teams() (res map[string]float64, err error) {
// 	res, err = d.singleLabelQuery(`
// 		SELECT
// 			count(*)
// 		FROM
// 			teams
// 		GROUP BY
// 			paused
// 	`)
// 	if err != nil {
// 		err = errors.Wrapf(err,
// 			"failed to pipelines by team")
// 		return
// 	}

// 	return
// }
