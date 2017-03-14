package main

import (
	"database/sql"
	"fmt"

	"github.com/nullbio/mig"
)

func Up(db *sql.DB, dir string) error {
	migrations, err := mig.CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}

	for {
		current, err := mig.GetDBVersion(db)
		if err != nil {
			return err
		}

		next, err := migrations.Next(current)
		if err != nil {
			if err == mig.ErrNoNextVersion {
				fmt.Printf("mig: no migrations to run. current version: %d\n", current)
				return nil
			}
			return err
		}

		if err = next.Up(db); err != nil {
			return err
		}
	}

	return nil
}

func UpByOne(db *sql.DB, dir string) error {
	migrations, err := mig.CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}

	currentVersion, err := mig.GetDBVersion(db)
	if err != nil {
		return err
	}

	next, err := migrations.Next(currentVersion)
	if err != nil {
		if err == mig.ErrNoNextVersion {
			fmt.Printf("mig: no migrations to run. current version: %d\n", currentVersion)
		}
		return err
	}

	if err = next.Up(db); err != nil {
		return err
	}

	return nil
}
