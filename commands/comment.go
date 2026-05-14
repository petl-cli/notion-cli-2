package commands

import "github.com/spf13/cobra"

var commentCmd = &cobra.Command{
	Use:   "comment",
	Short: "",
}

func init() {
	rootCmd.AddCommand(commentCmd)
}
