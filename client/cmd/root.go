package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

func init() {
	addFlag(rootCmd.Flags(), &flag{name: "server-url", short: "s", defaultValue: "http://localhost:9090", desc: "Filestore server url"})
	addFlag(rootCmd.Flags(), &flag{name: "log-level", short: "l", desc: "logging verbosity", defaultValue: "info"})
}
// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "store",
	Short: "store is a tool to operate filestore server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// Execute adds all child commands to the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}