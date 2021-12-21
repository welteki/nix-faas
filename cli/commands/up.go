package commands

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/spf13/cobra"
	"github.com/welteki/nix-faas/cli/nix"
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

	config, err := readNixFaasConfig(stackYaml)

	for _, image := range config.StackMetadata.Images {
		err := push(image)
		if err != nil {
			return err
		}
	}

	err = deploy(stackYaml)
	if err != nil {
		return err
	}

	return nil
}

func deploy(yamlFile string) error {
	cmd := "faas-cli"

	args := []string{
		"deploy",
		"-f",
		yamlFile,
	}

	task := execute.ExecTask{
		Command:     cmd,
		Args:        args,
		StreamStdio: true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("%q terminated with non-zero exit code", cmd)
	}

	return nil
}
