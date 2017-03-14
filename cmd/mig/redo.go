package main

import (
	"fmt"

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
	driver, conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	name, err := mig.Redo(driver, conn, viper.GetString("dir"))
	if err != nil {
		fmt.Printf("Success   %v\n", name)
	}

	return err
}
