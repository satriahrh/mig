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

var downAllCmd = &cobra.Command{
	Use:     "downall",
	Short:   "Roll back all migrations",
	Long:    "Roll back all migrations",
	Example: `mig downall sqlite3 ./foo.db`,
	RunE:    downRunE,
}

func init() {
	downCmd.Flags().StringP("dir", "d", ".", "directory with migration files")
	downAllCmd.Flags().StringP("dir", "d", ".", "directory with migration files")

	rootCmd.AddCommand(downCmd)
	rootCmd.AddCommand(downAllCmd)

	downCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(downCmd.Flags())
	}
	downAllCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(downAllCmd.Flags())
	}
}

func downRunE(cmd *cobra.Command, args []string) error {
	driver, conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	name, err := mig.Down(driver, conn, viper.GetString("dir"))
	if err != nil {
		fmt.Printf("Success   %v\n", name)
	}

	return err
}

func downAllRunE(cmd *cobra.Command, args []string) error {
	driver, conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	count, err := mig.DownAll(driver, conn, viper.GetString("dir"))
	if err != nil {
		fmt.Printf("Success   %d migrations\n", count)
	}

	return err
}
