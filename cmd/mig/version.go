package main

import (
	"fmt"

	"github.com/satriahrh/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var versionCmd = &cobra.Command{
	Use:     "version",
	Short:   "Print the current version of the database",
	Long:    "Print the current version of the database",
	Example: `$ mig version "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true"`,
	RunE:    versionRunE,
}

func init() {
	rootCmd.AddCommand(versionCmd)
	versionCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(versionCmd.Flags())
	}
}

func versionRunE(cmd *cobra.Command, args []string) error {
	conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	version, err := mig.Version(conn)
	if err != nil {
		return err
	}

	if version == 0 {
		fmt.Println("No migrations applied")
	} else {
		fmt.Printf("Version %d\n", version)
	}

	return nil
}
