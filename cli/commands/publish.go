package commands

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/spf13/cobra"
	"github.com/welteki/nix-faas/cli/image"
	"github.com/welteki/nix-faas/cli/nix"
	"github.com/welteki/nix-faas/cli/stack"
)

func init() {
	rootCmd.AddCommand(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:     "publish -f MODULE_FILE",
	Short:   "Builds and pushes OpenFaas function images to remote registry.",
	PreRunE: preRunPublish,
	RunE:    runPublish,
}

func preRunPublish(cmd *cobra.Command, args []string) error {
	if len(stackModule) == 0 {
		return fmt.Errorf("--file or -f is required")
	}

	return nil
}

func runPublish(cmd *cobra.Command, args []string) (retErr error) {
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

	return err
}

func push(m stack.ImageMetadata) (retErr error) {
	a, err := image.NewArchiveFromStream(m.Source)
	if err != nil {
		return fmt.Errorf("creating image archive: %w", err)
	}
	defer func() {
		if err := a.Close(); err != nil {
			retErr = fmt.Errorf("(archive: %v): %w", err, retErr)
		}
	}()

	cmd := "skopeo"

	args := []string{
		"copy",
		fmt.Sprintf("docker-archive:%s", a.Path()),
		fmt.Sprintf("docker://%s", m.Specifier),
		"--insecure-policy",
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
		return fmt.Errorf("received not-zero exit code from evaluation, error: %s", res.Stderr)
	}

	return nil
}

func readNixFaasConfig(path string) (stack.NixFaas, error) {
	config := stack.NixFaas{}

	configBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("reading file %q: %w", path, err)
	}
	unmarshallErr := yaml.Unmarshal(configBytes, &config)
	if unmarshallErr != nil {
		return config, fmt.Errorf("reading configuration: %w", err)
	}
	return config, nil
}
