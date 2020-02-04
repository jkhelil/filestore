package cmd

import (
	"filestore/client/store"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(RegisterFrequentCommand())
}

// RegisterFrequentCommand register frequent subcommand and flags 
func RegisterFrequentCommand() *cobra.Command {
	c := &cobra.Command{
		Use:  "freq-words",
		Run: func(cmd *cobra.Command, args []string) {
			c := store.NewClient()
			c.FreqWords()
		},
	}
	addFlag(c.Flags(), &flag{name: "limit", short: "n", desc: "limit for frequent words", defaultValue: 1, kind: "int"})
	addFlag(c.Flags(), &flag{name: "order", desc: "order for frequent words", defaultValue: "dsc"})
	return c
}
