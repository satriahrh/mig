package mig

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

var (
	ErrNoCurrentVersion = errors.New("no current version found")
	ErrNoNextVersion    = errors.New("no next version found")
)

var Log = os.Stdout

type Migrations []*Migration

// helpers so we can use pkg sort
func (ms Migrations) Len() int      { return len(ms) }
func (ms Migrations) Swap(i, j int) { ms[i], ms[j] = ms[j], ms[i] }
func (ms Migrations) Less(i, j int) bool {
	if ms[i].Version == ms[j].Version {
		panic(fmt.Sprintf("mig: duplicate version %v detected:\n%v\n%v", ms[i].Version, ms[i].Source, ms[j].Source))
	}
	return ms[i].Version < ms[j].Version
}

func (ms Migrations) Current(current int64) (*Migration, error) {
	for i, migration := range ms {
		if migration.Version == current {
			return ms[i], nil
		}
	}

	return nil, ErrNoCurrentVersion
}

func (ms Migrations) Next(current int64) (*Migration, error) {
	for i, migration := range ms {
		if migration.Version > current {
			return ms[i], nil
		}
	}

	return nil, ErrNoNextVersion
}

func (ms Migrations) Last() (*Migration, error) {
	if len(ms) == 0 {
		return nil, ErrNoNextVersion
	}

	return ms[len(ms)-1], nil
}

func (ms Migrations) String() string {
	str := ""
	for _, m := range ms {
		str += fmt.Sprintln(m)
	}
	return str
}

// collect all the valid looking migration scripts in the migrations folder,
// and order them by version.
func collectMigrations(dirpath string, current, target int64) (Migrations, error) {
	var migrations Migrations

	// extract the numeric component of each migration,
	// filter out any uninteresting files,
	// and ensure we only have one file per migration version.
	files, err := filepath.Glob(dirpath + "/*.sql")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		v, err := NumericComponent(file)
		if err != nil {
			return nil, err
		}
		if versionFilter(v, current, target) {
			migration := &Migration{Version: v, Next: -1, Previous: -1, Source: file}
			migrations = append(migrations, migration)
		}
	}

	migrations = sortAndConnectMigrations(migrations)

	return migrations, nil
}

// sortAndConnectMigrations sorts the migrations based on the version numbers
// and creates a linked list between each migration.
func sortAndConnectMigrations(migrations Migrations) Migrations {
	// Sort the migrations based on version
	sort.Sort(migrations)

	// now that we're sorted in the appropriate direction,
	// populate next and previous for each migration
	for i, m := range migrations {
		prev := int64(-1)
		if i > 0 {
			prev = migrations[i-1].Version
			migrations[i-1].Next = m.Version
		}
		migrations[i].Previous = prev
	}

	return migrations
}

// versionFilter returns true if v is within the current version and target
// version range.
func versionFilter(v, current, target int64) bool {
	if target > current {
		return v > current && v <= target
	}

	if target < current {
		return v <= current && v > target
	}

	return false
}

// Create the mig_migrations table
// and insert the initial 0 value into it
func createVersionTable(db *sql.DB) error {
	txn, err := db.Begin()
	if err != nil {
		return err
	}

	d := getDialect()

	if _, err := txn.Exec(d.createVersionTableSQL()); err != nil {
		txn.Rollback()
		return err
	}

	version := 0
	applied := true
	if _, err := txn.Exec(d.insertVersionSQL(), version, applied); err != nil {
		txn.Rollback()
		return err
	}

	return txn.Commit()
}

// getVersion retrieves the current version for this database.
// Create and initialize the database migration table if it doesn't exist.
func getVersion(db *sql.DB) (int64, error) {
	rows, err := getDialect().versionQuery(db)
	if err != nil {
		return 0, createVersionTable(db)
	}
	defer rows.Close()

	// The most recent record for each migration specifies
	// whether it has been applied or rolled back.
	// The first version we find that has been applied is the current version.

	toSkip := make([]int64, 0)

	for rows.Next() {
		var row MigrationRecord
		if err = rows.Scan(&row.VersionId, &row.IsApplied); err != nil {
			return 0, errors.New("error scanning rows:", err)
		}

		// have we already marked this version to be skipped?
		skip := false
		for _, v := range toSkip {
			if v == row.VersionId {
				skip = true
				break
			}
		}

		if skip {
			continue
		}

		// if version has been applied we're done
		if row.IsApplied {
			return row.VersionId, nil
		}

		// latest version of migration has not been applied.
		toSkip = append(toSkip, row.VersionId)
	}

	panic("unreachable")
}
