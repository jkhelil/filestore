package cmd

import (
	"os"

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
			if err := c.List(); err != nil {
				os.Exit(1)
			}
		},
	}

	return c
}
