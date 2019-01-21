package main

import (
	"fmt"

	"github.com/satriahrh/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Roll back the version by one",
	Long:  "Roll back the version by one",
	Example: `$ mig down "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
"`,
	RunE: downRunE,
}

var downAllCmd = &cobra.Command{
	Use:   "downall",
	Short: "Roll back all migrations",
	Long:  "Roll back all migrations",
	Example: `$ mig downall "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
"`,
	RunE: downAllRunE,
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
	conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	name, err := mig.Down(conn, viper.GetString("dir"))
	if mig.IsNoMigrationError(err) {
		fmt.Println("No migrations to run")
		return nil
	} else if err != nil {
		return err
	}

	fmt.Printf("Success   %v\n", name)
	return nil
}

func downAllRunE(cmd *cobra.Command, args []string) error {
	conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	count, err := mig.DownAll(conn, viper.GetString("dir"))
	if err != nil {
		return err
	}

	if count == 0 {
		fmt.Println("No migrations to run")
	} else {
		fmt.Printf("Success   %d migrations\n", count)
	}

	return nil
}
