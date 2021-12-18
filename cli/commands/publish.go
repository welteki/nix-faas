package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"gopkg.in/yaml.v2"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/spf13/cobra"
	"github.com/welteki/nix-faas/cli/image"
	"github.com/welteki/nix-faas/cli/stack"
)

func init() {
	build, _, _ := rootCmd.Find([]string{"build"})
	publishCmd.Flags().AddFlagSet(build.Flags())

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

func runPublish(cmd *cobra.Command, args []string) error {
	var out string
	if len(outLink) == 0 {
		// TODO: outlink should contain random part
		out = path.Join(os.TempDir(), ".tmp-stack.yaml")
		defer os.Remove(outLink)
	} else {
		out = outLink
	}

	err := build(stackModule, out, false)
	if err != nil {
		return err
	}

	config, err := readNixFaasConfig(out)

	for _, image := range config.StackMetadata.Images {
		err := push(image)
		if err != nil {
			return err
		}
	}

	return err
}

func push(m stack.ImageMetadata) error {
	a, err := image.NewArchiveFromStream(m.Source)

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

func readNixFaasConfig(yamlFile string) (stack.NixFaas, error) {
	config := stack.NixFaas{}

	configBytes, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return config, fmt.Errorf("Error while reading files %s", err.Error())
	}
	unmarshallErr := yaml.Unmarshal(configBytes, &config)
	if unmarshallErr != nil {
		return config, fmt.Errorf("Error while reading configuration: %s", err.Error())
	}
	return config, nil
}
