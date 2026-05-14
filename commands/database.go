package commands

import "github.com/spf13/cobra"

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "",
}

func init() {
	rootCmd.AddCommand(databaseCmd)
}
