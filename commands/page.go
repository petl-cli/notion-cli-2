package commands

import "github.com/spf13/cobra"

var pageCmd = &cobra.Command{
	Use:   "page",
	Short: "",
}

func init() {
	rootCmd.AddCommand(pageCmd)
}
