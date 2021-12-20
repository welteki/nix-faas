package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/welteki/nix-faas/cli/nix"
)

func init() {
	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:     "build -f MODULE_FILE",
	Short:   "Build OpenFaaS functions",
	PreRunE: preRunBuild,
	RunE:    runBuild,
}

func preRunBuild(cmd *cobra.Command, args []string) error {
	if len(stackModule) == 0 {
		return fmt.Errorf("--file or -f is required")
	}

	return nil
}

func runBuild(cmd *cobra.Command, args []string) (retErr error) {
	err := nix.BuildStack(stackModule, nil)

	if err != nil {
		return err
	}

	return nil
}
