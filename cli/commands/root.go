package commands

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	stackModules []string
)

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&stackModules, "file", "f", []string{}, "Path to nix module file describing function(s)")
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "nix-faas",
	Short: "Manage your OpenFaaS functions using nix",
	Long: `
Manage your OpenFaaS functions using nix`,
	Run: runRootCmd,
}

func runRootCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}
