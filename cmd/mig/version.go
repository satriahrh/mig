package main

import (
	"database/sql"
	"fmt"

	"github.com/nullbio/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the current version of the database",
	Long:    "Print the current version of the database",
	Example: `mig version postgres "user=postgres dbname=postgres sslmode=disable"`,
	RunE:    versionRunE,
}

func init() {
	rootCmd.AddCommand(upCmd)
	versionCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(versionCmd.Flags())
	}
}

func versionRunE(cmd *cobra.Command, args []string) error {
	return nil
}

func Version(db *sql.DB, dir string) error {
	current, err := mig.GetDBVersion(db)
	if err != nil {
		return err
	}

	fmt.Printf("mig: version %v\n", current)
	return nil
}
