package cmd

import (
	"os"

	"filestore/client/store"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(RegisterRemoveCommand())
}

// RegisterRemoveCommand register Remove subcommand and flags 
func RegisterRemoveCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "rm",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient()	
			if err := c.Remove(args[0]); err != nil {
				os.Exit(1)
			}
		},
	}

	return c
}
