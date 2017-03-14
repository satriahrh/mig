package main

import (
	"database/sql"
	"fmt"

	"github.com/nullbio/mig"
)

func Version(db *sql.DB, dir string) error {
	current, err := mig.GetDBVersion(db)
	if err != nil {
		return err
	}

	fmt.Printf("mig: version %v\n", current)
	return nil
}
