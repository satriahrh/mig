package main

import (
	"fmt"

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
	driver, conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	status, err := mig.Status(driver, conn)
	if err != nil {
		return err
	}

	fmt.Println("    Applied At                  Migration")
	fmt.Println("    =======================================")
	for _, s := range status {
		fmt.Printf("%-24s -- %v\n", s.Applied, s.Name)
	}

	return nil
}
