package commands

import "github.com/spf13/cobra"

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "",
}

func init() {
	rootCmd.AddCommand(tokenCmd)
}
