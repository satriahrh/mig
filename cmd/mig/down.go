package main

import (
	"fmt"

	"github.com/nullbio/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downCmd = &cobra.Command{
	Use:     "down",
	Short:   "Roll back the version by one",
	Long:    "Roll back the version by one",
	Example: `mig down sqlite3 ./foo.db`,
	RunE:    downRunE,
}

func init() {
	downCmd.Flags().StringP("dir", "d", ".", "directory with migration files")

	rootCmd.AddCommand(downCmd)
	downCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(downCmd.Flags())
	}
}

func downRunE(cmd *cobra.Command, args []string) error {
	driver, conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	name, err := mig.Down(driver, conn, viper.GetString("dir"))
	if err != nil {
		fmt.Printf("Success %v\n", name)
	}

	return err
}
