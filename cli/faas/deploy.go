package faas

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func Deploy(yamlFile string) error {
	cmd, args := getDeployCommand(yamlFile)

	task := execute.ExecTask{
		Command:     cmd,
		Args:        args,
		StreamStdio: true,
	}

	res, err := task.Execute()

	if err != nil {
		return fmt.Errorf("executing %q: %w", cmd, err)
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("%q terminated with non-zero exit code", cmd)
	}

	return nil
}

func getDeployCommand(yamlFile string) (string, []string) {
	args := []string{
		"deploy",
		"-f",
		yamlFile,
	}

	cmd := "faas-cli"

	return cmd, args
}
