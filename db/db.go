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
