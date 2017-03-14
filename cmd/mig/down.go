package main

import (
	"database/sql"
	"fmt"

	"github.com/nullbio/mig"
)

func Down(db *sql.DB, dir string) error {
	currentVersion, err := mig.GetDBVersion(db)
	if err != nil {
		return err
	}

	migrations, err := mig.CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}

	current, err := migrations.Current(currentVersion)
	if err != nil {
		return fmt.Errorf("no migration %v", currentVersion)
	}

	return current.Down(db)
}
