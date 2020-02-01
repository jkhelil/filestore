package cmd

import (
	"filestore/client/store"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func init() {
	rootCmd.AddCommand(RegisterListCommand())
}

// RegisterListCommand register add subcommand and flags 
func RegisterListCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "ls",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient(viper.GetString("server-url"))
			c.List()
		},
	}

	return c
}
