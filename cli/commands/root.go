package commands

import (
	"os"

	"github.com/spf13/cobra"
)

// Path to nix files
var NixDir string

var (
	stackModule string
)

func init() {
	rootCmd.PersistentFlags().StringVarP(&stackModule, "file", "f", "", "Path to nix module file describing function(s)")
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

func getNixDir() string {
	if len(NixDir) == 0 {
		return os.Getenv("NIXFAAS_NIX_DIR")
	}

	return NixDir
}
