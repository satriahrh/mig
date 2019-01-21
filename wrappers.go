package mig

import (
	"database/sql"
	"fmt"
	"math"
	"path/filepath"
	"time"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

// Create a templated migration file in dir
func Create(name, dir string) (string, error) {
	timestamp := time.Now().Format("20060102150405")
	filename := fmt.Sprintf("%v_%v.sql", timestamp, name)

	fpath := filepath.Join(dir, filename)
	tmpl := migrationTemplate

	path, err := writeTemplateToFile(fpath, tmpl, timestamp)
	return path, err
}

// Down rolls back the version by one
func Down(conn, dir string) (name string, err error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return "", err
	}

	err = setDialect()
	if err != nil {
		return "", err
	}

	return DownDB(db, dir)
}

// DownDB rolls back the version by one
// Expects SetDialect to be called beforehand.
func DownDB(db *sql.DB, dir string) (name string, err error) {
	currentVersion, err := getVersion(db)
	if err != nil {
		return "", err
	}

	migrations, err := collectMigrations(dir, 0, math.MaxInt64)
	if err != nil {
		return "", err
	}

	current, err := migrations.current(currentVersion)
	if err != nil {
		return "", errNoMigration{}
	}

	return current.down(db)
}

// DownAll rolls back all migrations.
// Logs success messages to global writer variable Log.
func DownAll(conn, dir string) (int, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return 0, err
	}

	err = setDialect()
	if err != nil {
		return 0, err
	}

	return DownAllDB(db, dir)
}

// DownAllDB rolls back all migrations.
// Logs success messages to global writer variable Log.
// Expects SetDialect to be called beforehand.
func DownAllDB(db *sql.DB, dir string) (int, error) {
	count := 0

	migrations, err := collectMigrations(dir, 0, math.MaxInt64)
	if err != nil {
		return count, err
	}

	for {
		currentVersion, err := getVersion(db)
		if err != nil {
			return count, err
		}

		current, err := migrations.current(currentVersion)
		// no migrations left to run
		if err != nil {
			return count, nil
		}

		name, err := current.down(db)
		if err != nil {
			return count, err
		}

		Log.Write([]byte(fmt.Sprintf("Success   %v\n", name)))
		count++
	}
}

// Up migrates to the highest version available
func Up(conn, dir string) (int, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return 0, err
	}

	err = setDialect()
	if err != nil {
		return 0, err
	}

	return UpDB(db, dir)
}

// UpDB migrates to the highest version available
// Expects SetDialect to be called beforehand.
func UpDB(db *sql.DB, dir string) (int, error) {
	count := 0

	migrations, err := collectMigrations(dir, 0, math.MaxInt64)
	if err != nil {
		return count, err
	}

	for {
		currentVersion, err := getVersion(db)
		if err != nil {
			return count, err
		}

		next, err := migrations.next(currentVersion)
		// no migrations left to run
		if err != nil {
			return count, nil
		}

		name, err := next.up(db)
		if err != nil {
			return count, err
		}

		Log.Write([]byte(fmt.Sprintf("Success   %v\n", name)))
		count++
	}
}

// UpOne migrates one version
func UpOne(conn, dir string) (name string, err error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return "", err
	}

	err = setDialect()
	if err != nil {
		return "", err
	}

	return UpOneDB(db, dir)
}

// UpOneDB migrates one version
// Expects SetDialect to be called beforehand.
func UpOneDB(db *sql.DB, dir string) (name string, err error) {
	currentVersion, err := getVersion(db)
	if err != nil {
		return "", err
	}

	migrations, err := collectMigrations(dir, 0, math.MaxInt64)
	if err != nil {
		return "", err
	}

	next, err := migrations.next(currentVersion)
	if err != nil {
		return "", errNoMigration{}
	}

	return next.up(db)
}

// Redo re-runs the latest migration.
func Redo(conn, dir string) (string, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return "", err
	}

	err = setDialect()
	if err != nil {
		return "", err
	}

	return RedoDB(db, dir)
}

// RedoDB re-runs the latest migration.
// Expects SetDialect to be called beforehand.
func RedoDB(db *sql.DB, dir string) (string, error) {
	currentVersion, err := getVersion(db)
	if err != nil {
		return "", err
	}

	migrations, err := collectMigrations(dir, 0, math.MaxInt64)
	if err != nil {
		return "", err
	}

	current, err := migrations.current(currentVersion)
	if err != nil {
		return "", errNoMigration{}
	}

	if _, err := current.down(db); err != nil {
		return "", err
	}

	return current.up(db)
}

// MigrationStatus show the status of the migration
type MigrationStatus struct {
	Applied string
	Name    string
}

// Status returns the status of each migration
func Status(conn, dir string) ([]MigrationStatus, error) {
	s := []MigrationStatus{}

	db, err := sql.Open("mysql", conn)
	if err != nil {
		return s, err
	}

	err = setDialect()
	if err != nil {
		return s, err
	}

	return StatusDB(db, dir)
}

// StatusDB returns the status of each migration
// Expects SetDialect to be called beforehand
func StatusDB(db *sql.DB, dir string) ([]MigrationStatus, error) {
	s := []MigrationStatus{}

	migrations, err := collectMigrations(dir, 0, math.MaxInt64)
	if err != nil {
		return s, err
	}

	// must ensure that the version table exists if we're running on a pristine DB
	if _, err := getVersion(db); err != nil {
		return s, err
	}

	for _, migration := range migrations {
		s = append(s, MigrationStatus{
			Applied: getMigrationStatus(db, migration.version),
			Name:    filepath.Base(migration.source),
		})
	}

	return s, nil
}

// Version returns the current migration version
func Version(conn string) (int64, error) {
	db, err := sql.Open("mysql", conn)
	if err != nil {
		return 0, err
	}

	err = setDialect()
	if err != nil {
		return 0, err
	}

	return VersionDB(db)
}

// VersionDB returns the current migration version
// Expects SetDialect to be called beforehand
func VersionDB(db *sql.DB) (int64, error) {
	return getVersion(db)
}
