package commands

import "github.com/spf13/cobra"

var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "",
}

func init() {
	rootCmd.AddCommand(searchCmd)
}
