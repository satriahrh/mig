package main

import (
	"database/sql"
	"fmt"

	"github.com/nullbio/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var upCmd = &cobra.Command{
	Use:     "up",
	Short:   "Migrate the database to the most recent version available",
	Long:    "Migrate the database to the most recent version available",
	Example: `mig up mysql "user:password@/dbname"`,
	RunE:    upRunE,
}

func init() {
	upCmd.Flags().StringP("dir", "d", ".", "directory with migration files")

	rootCmd.AddCommand(upCmd)
	upCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(upCmd.Flags())
	}
}

func upRunE(cmd *cobra.Command, args []string) error {
	return nil
}

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
