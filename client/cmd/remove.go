package cmd

import (
	"filestore/client/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(RegisterRemoveCommand())
}

// RegisterRemoveCommand register add subcommand and flags 
func RegisterRemoveCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "rm",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient(viper.GetString("server-url"))
			c.Remove(args[0])
		},
	}

	return c
}
