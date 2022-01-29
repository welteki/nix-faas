package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	// GitCommit Git Commit SHA
	GitCommit string
	// Version version of the CLI
	Version string
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version information.",
	Run:   runVersionCommand,
}

func runVersionCommand(cmd *cobra.Command, args []string) {
	printVersion()
}

func printVersion() {
	fmt.Printf("nix-faas version: %s\tcommit: %s\n", getVersion(), GitCommit)
}

func getVersion() string {
	if len(Version) == 0 {
		return "dev"
	}
	return Version
}
