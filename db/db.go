package db

import (
	"database/sql"

	"github.com/pkg/errors"

	_ "github.com/lib/pq"
)

type Db struct {
	db *sql.DB
}

func New(connStr string) (db *Db, err error) {
	dbConn, err := sql.Open("postgres", connStr)
	if err != nil {
		err = errors.Wrapf(err,
			"couldn't open db")
		return
	}

	db = &Db{db: dbConn}

	return
}

func (d *Db) Close() (err error) {
	err = d.db.Close()
	if err != nil {
		err = errors.Wrapf(err,
			"failed to close database")
		return
	}

	return
}

type Datapoint struct {
	LabelSet []string
	Value    float64
}

func (d *Db) query(query string, labelCount int) (res []Datapoint, err error) {
	rows, err := d.db.Query(query)
	if err != nil {
		err = errors.Wrapf(err,
			"failed to execute query")
		return
	}

	var (
		count           float64 = 0
		scanningTargets         = []interface{}{&count}
		labelSet                = make([]string, labelCount)
	)

	for i := 0; i < labelCount; i++ {
		scanningTargets = append(scanningTargets, &labelSet[i])
	}

	for rows.Next() {
		err = rows.Scan(scanningTargets...)
		if err != nil {
			err = errors.Wrapf(err,
				"failed to scan query row")
			return
		}

		dp := Datapoint{Value: count}
		dp.LabelSet = make([]string, labelCount)
		copy(dp.LabelSet, labelSet)

		res = append(res, dp)
	}

	return
}
