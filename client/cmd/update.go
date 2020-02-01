package cmd

import (
	"filestore/client/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(RegisterUpdateCommand())
}

// RegisterUpdateommand register add subcommand and flags 
func RegisterUpdateCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "update",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient(viper.GetString("server-url"))
			c.Update(args[0])
		},
	}
	return c
}
