package main

import (
	"database/sql"
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/nullbio/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var statusCmd = &cobra.Command{
	Use:     "status",
	Short:   "Dump the migration status for the database",
	Long:    "Dump the migration status for the database",
	Example: `mig status postgres "user=postgres dbname=postgres sslmode=disable"`,
	RunE:    statusRunE,
}

func init() {
	rootCmd.AddCommand(statusCmd)
	statusCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(statusCmd.Flags())
	}
}

func statusRunE(cmd *cobra.Command, args []string) error {
	return nil
}

func Status(db *sql.DB, dir string) error {
	// Collect all migrations
	migrations, err := mig.CollectMigrations(dir, minVersion, maxVersion)
	if err != nil {
		return err
	}

	// must ensure that the version table exists if we're running on a pristine DB
	if _, err := getVersion(db); err != nil {
		return err
	}

	fmt.Println("    Applied At                  Migration")
	fmt.Println("    =======================================")
	for _, migration := range migrations {
		printMigrationStatus(db, migration.Version, filepath.Base(migration.Source))
	}

	return nil
}

func printMigrationStatus(db *sql.DB, version int64, script string) {
	var row mig.MigrationRecord
	q := fmt.Sprintf("SELECT tstamp, is_applied FROM mig_migrations WHERE version_id=%d ORDER BY tstamp DESC LIMIT 1", version)
	e := db.QueryRow(q).Scan(&row.TStamp, &row.IsApplied)

	if e != nil && e != sql.ErrNoRows {
		log.Fatal(e)
	}

	var appliedAt string

	if row.IsApplied {
		appliedAt = row.TStamp.Format(time.ANSIC)
	} else {
		appliedAt = "Pending"
	}

	fmt.Printf("    %-24s -- %v\n", appliedAt, script)
}
