package mig

import (
	"database/sql"
)

// sqlDialect abstracts the details of specific SQL dialects
// for mig's few SQL specific statements
type sqlDialect interface {
	createVersionTableSQL() string // sql string to create the mig_migrations table
	insertVersionSQL() string      // sql string to insert the initial version table row
	versionQuery(db *sql.DB) (*sql.Rows, error)
}

var dialect sqlDialect = &mySQLDialect{}

func getDialect() sqlDialect {
	return dialect
}

// SetDialect sets the current driver dialect for all future calls
// to the library.
func SetDialect() error {
	dialect = &mySQLDialect{}
	return nil
}

// SetDialect sets the current driver dialect for all future calls
// to the library.
func setDialect() error {
	dialect = &mySQLDialect{}
	return nil
}

type mySQLDialect struct{}

func (mySQLDialect) createVersionTableSQL() string {
	return `CREATE TABLE mig_migrations (
                id serial NOT NULL,
                version_id bigint NOT NULL,
                is_applied boolean NOT NULL,
                tstamp timestamp NULL default now(),
                PRIMARY KEY(id)
            );`
}

func (mySQLDialect) insertVersionSQL() string {
	return "INSERT INTO mig_migrations (version_id, is_applied) VALUES (?, ?);"
}

func (mySQLDialect) versionQuery(db *sql.DB) (*sql.Rows, error) {
	rows, err := db.Query("SELECT version_id, is_applied from mig_migrations ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	return rows, err
}
