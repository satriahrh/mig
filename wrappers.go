package mig

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/spf13/viper"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
	// postgres driver
	_ "github.com/lib/pq"
	// sqlite3 driver
	_ "github.com/mattn/go-sqlite3"
)

type errNoMigration struct {
	version int64
}

func (e errNoMigration) Error() string {
	return fmt.Sprintf("no migration %d", e.version)
}

func IsNoMigrationError(err error) bool {
	_, ok := err.(errNoMigration)
	if ok {
		return ok
	}

	return false
}

// Down rolls back the version by one
func Down(driver string, conn string, dir string) (name string, err error) {
	db, err := sql.Open(driver, conn)
	if err != nil {
		return "", err
	}

	err = setDialect(driver)
	if err != nil {
		return "", err
	}

	currentVersion, err := getVersion(db)
	if err != nil {
		return "", err
	}

	migrations, err := collectMigrations(viper.GetString("dir"), 0, math.MaxInt64)
	if err != nil {
		return "", err
	}

	current, err := migrations.Current(currentVersion)
	if err != nil {
		return "", errNoMigration{version: current.Version}
	}

	return current.Down(db)
}
