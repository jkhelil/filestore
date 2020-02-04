package cmd

import (
	"filestore/client/store"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(RegisterCountCommand())
}

// RegisterCountCommand register count subcommand and flags 
func RegisterCountCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "wc",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient()
			c.CountWords()
		},
	}
	return c
}
