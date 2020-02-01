package cmd

import (
	"filestore/client/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(RegisterAddCommand())
}

// RegisterAddCommand register add subcommand and flags 
func RegisterAddCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "add",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient(viper.GetString("server-url"))
			c.Add(args)
		},
	}
	return c
}
