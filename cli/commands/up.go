package commands

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/welteki/nix-faas/cli/faas"
	"github.com/welteki/nix-faas/cli/nix"
	"github.com/welteki/nix-faas/cli/stack"
)

func init() {
	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:     "up -f MODULE_FILE",
	Short:   "Build, push and deploy OpenFaaS functions",
	PreRunE: preRunUp,
	RunE:    runUp,
}

func preRunUp(cmd *cobra.Command, args []string) error {
	if len(stackModule) == 0 {
		return fmt.Errorf("--file or -f is required")
	}

	return nil
}

func runUp(cmd *cobra.Command, args []string) (retErr error) {
	gcRoot, err := nix.NewGarbageCollectionRoot()
	if err != nil {
		return err
	}
	defer func() {
		if err := gcRoot.Close(); err != nil {
			retErr = fmt.Errorf("(gcroot: %v): %w", err, retErr)
		}
	}()

	err = nix.BuildStack(stackModule, gcRoot)
	if err != nil {
		return err
	}

	stackYaml := gcRoot.Path()

	config, err := stack.ReadNixFaasConfig(stackYaml)
	if err != nil {
		return fmt.Errorf("getting nix-faas config: %w", err)
	}

	for _, image := range config.StackMetadata.Images {
		err := push(image)
		if err != nil {
			return err
		}
	}

	err = faas.Deploy(stackYaml)
	if err != nil {
		return fmt.Errorf("deploying functions: %w", err)
	}

	return nil
}
