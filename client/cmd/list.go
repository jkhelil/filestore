package cmd

import (
	"filestore/client/store"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(RegisterListCommand())
}

// RegisterListCommand register List subcommand and flags 
func RegisterListCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "ls",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient()
			c.List()
		},
	}

	return c
}
