package main

import (
	"fmt"

	"github.com/satriahrh/mig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var redoCmd = &cobra.Command{
	Use:   "redo",
	Short: "Down then up the latest migration",
	Long:  "Down then up the latest migration",
	Example: `$ mig redo "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
"`,
	RunE: redoRunE,
}

var redoAllCmd = &cobra.Command{
	Use:   "redo",
	Short: "Down then up all migrations",
	Long:  "Down then up all migrations",
	Example: `$ mig redoall "user:password@tcp(localhost:5555)/dbname?tls=skip-verify&autocommit=true
"`,
	RunE: redoAllRunE,
}

func init() {
	redoCmd.Flags().StringP("dir", "d", ".", "directory with migration files")

	rootCmd.AddCommand(redoCmd)
	redoCmd.PreRun = func(*cobra.Command, []string) {
		viper.BindPFlags(redoCmd.Flags())
	}
}

func redoRunE(cmd *cobra.Command, args []string) error {
	conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	name, err := mig.Redo(conn, viper.GetString("dir"))
	if mig.IsNoMigrationError(err) {
		fmt.Println("No migrations to run")
		return nil
	} else if err != nil {
		return err
	}

	fmt.Printf("Success   %v\n", name)
	return nil
}

func redoAllRunE(cmd *cobra.Command, args []string) error {
	conn, err := getConnArgs(args)
	if err != nil {
		return err
	}

	_, err = mig.DownAll(conn, viper.GetString("dir"))
	if err != nil {
		return err
	}

	count, err := mig.Up(conn, viper.GetString("dir"))
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
