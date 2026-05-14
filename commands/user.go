package commands

import "github.com/spf13/cobra"

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "",
}

func init() {
	rootCmd.AddCommand(userCmd)
}
