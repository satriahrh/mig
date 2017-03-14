package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/nullbio/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var createCmd = &cobra.Command{
	Use:     "create",
	Short:   "Create a blank migration template",
	Long:    "Create a blank migration template",
	Example: `mig create add_users`,
	RunE:    createRunE,
}

func init() {
	createCmd.Flags().StringP("dir", "d", ".", "directory with migration files")

	rootCmd.AddCommand(createCmd)
	createCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(createCmd.Flags())
	}
}

func createRunE(cmd *cobra.Command, args []string) error {
	if len(args) < 1 || len(args[0]) == 0 {
		return errors.New("no migration name provided")
	}

	path, err := mig.CreateMigration(args[0], viper.GetString("dir"), time.Now())
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Created migration %s", path))

	return nil
}
