package main

import (
	"database/sql"

	"github.com/nullbio/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var redoCmd = &cobra.Command{
	Use:     "redo",
	Short:   "Re-run the latest migration",
	Long:    "Re-run the latest migration",
	Example: `mig redo postgres "user=postgres dbname=postgres sslmode=disable"`,
	RunE:    redoRunE,
}

func init() {
	redoCmd.Flags().StringP("dir", "d", ".", "directory with migration files")

	rootCmd.AddCommand(redoCmd)
	redoCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(redoCmd.Flags())
	}
}

func redoRunE(cmd *cobra.Command, args []string) error {
	return nil
}

func Redo(db *sql.DB, dir string) error {
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
		return err
	}

	previous, err := migrations.Next(currentVersion)
	if err != nil {
		return err
	}

	if err := previous.Up(db); err != nil {
		return err
	}

	if err := current.Up(db); err != nil {
		return err
	}

	return nil
}
