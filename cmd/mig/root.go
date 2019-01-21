package main

import (
	"errors"
	"math"
	"os"

	"github.com/satriahrh/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	minVersion = int64(0)
	maxVersion = int64(math.MaxInt64)
)

var rootCmd = &cobra.Command{
	Use:   "mig",
	Short: "mig is a database migration tool for MySQL.",
	Long:  "mig is a database migration tool for MySQL.",
	Example: `$ mig up user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
$ mig down "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true"
$ mig create add_users`,
}

func init() {
	// Set the mig library logger to os.Stdout
	mig.Log = os.Stdout

	rootCmd.Flags().BoolP("version", "", false, "Print the mig tool version")
	viper.BindPFlags(rootCmd.Flags())
}

// getConnArgs takes in args from cobra and returns the 0th and 1st arg
// which should be the driver and connection string
func getConnArgs(args []string) (conn string, err error) {
	if len(args) < 1 {
		err = errors.New("no connection details provided")
		return
	}
	conn = args[0]

	return
}
