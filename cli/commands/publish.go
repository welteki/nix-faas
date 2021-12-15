package commands

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"

	"gopkg.in/yaml.v2"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/spf13/cobra"
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

func runPublish(cmd *cobra.Command, args []string) error {
	outLink := path.Join(os.TempDir(), ".tmp-stack.yaml")
	defer os.Remove(outLink)

	err := build(stackModule, outLink, false)
	if err != nil {
		return err
	}

	config, err := readNixFaasConfig(outLink)

	for _, image := range config.StackMetadata.Images {
		err := push(image)
		if err != nil {
			return err
		}
	}

	return err
}

func push(image stack.ImageMetadata) error {
	tarImageFile, err := ioutil.TempFile("", "nix-faas-docker-tar-*")
	if err != nil {
		return fmt.Errorf("Error while creating temporary file %s", err.Error())
	}
	defer os.Remove(tarImageFile.Name())

	generateImage := exec.Command(image.Source)
	generateImage.Stdout = tarImageFile

	startErr := generateImage.Start()
	if startErr != nil {
		return startErr
	}

	execErr := generateImage.Wait()
	if execErr != nil {
		return execErr
	}

	cmd := "skopeo"

	args := []string{
		"copy",
		fmt.Sprintf("docker-archive:%s", tarImageFile.Name()),
		fmt.Sprintf("docker://%s", image.Specifier),
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
