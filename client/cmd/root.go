package cmd

import (
	"os"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "store",
	Short: "store is a tool to operate filestore server",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

// registerStoreFlags register flags for store command
func registerStoreFlags(store *cobra.Command) {
	addFlag(store.PersistentFlags(), &flag{name: "server-url", defaultValue: "http://localhost:9090", desc: "Filestore server url"})
	addFlag(store.PersistentFlags(), &flag{name: "log-level", desc: "logging verbosity", defaultValue: "info"})
}

// Execute adds all child commands to the root command
func Execute() {
	registerStoreFlags(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}