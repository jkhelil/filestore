package cmd

import (
	"os"

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
			if err := c.CountWords(); err != nil {
				os.Exit(1)
			}
		},
	}
	return c
}
