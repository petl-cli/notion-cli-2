package commands

import "github.com/spf13/cobra"

var blockCmd = &cobra.Command{
	Use:   "block",
	Short: "",
}

func init() {
	rootCmd.AddCommand(blockCmd)
}
