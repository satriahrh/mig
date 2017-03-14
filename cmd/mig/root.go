package main

import (
	"errors"
	"math"
	"os"

	"github.com/nullbio/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	minVersion = int64(0)
	maxVersion = int64(math.MaxInt64)
)

var rootCmd = &cobra.Command{
	Use:   "mig",
	Short: "mig is a database migration tool for Postgres, MySQL and SQLite3.",
	Long:  "mig is a database migration tool for Postgres, MySQL and SQLite3.",
	Example: `mig postgres "user=postgres dbname=postgres sslmode=disable" up
mig mysql "user:password@/dbname" down
mig sqlite3 ./foo.db status
mig create add_users`,
}

func init() {
	// Set the mig library logger to os.Stdout
	mig.Log = os.Stdout

	rootCmd.Flags().BoolP("version", "", false, "Print the mig tool version")
	viper.BindPFlags(rootCmd.Flags())
}

// getConnArgs takes in args from cobra and returns the 0th and 1st arg
// which should be the driver and connection string
func getConnArgs(args []string) (driver string, conn string, err error) {
	if len(args) < 2 {
		return "", "", errors.New("no connection details provided")
	}

	return args[0], args[1], nil
}
