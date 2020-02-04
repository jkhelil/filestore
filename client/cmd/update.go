package cmd

import (
	"os"

	"filestore/client/store"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(RegisterUpdateCommand())
}

// RegisterUpdateCommand register update subcommand and flags 
func RegisterUpdateCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "update",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient()
			if err := c.Update(args[0]); err != nil {
				os.Exit(1)
			}
		},
	}
	return c
}
