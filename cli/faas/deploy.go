package faas

import (
	"fmt"
	"time"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func Deploy(yamlFile string, gateway string, timeout time.Duration, tlsInsecure bool) error {
	cmd, args := getDeployCommand(yamlFile, gateway, timeout, tlsInsecure)

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

func getDeployCommand(yamlFile string, gateway string, timeout time.Duration, tlsInsecure bool) (string, []string) {
	args := []string{
		"deploy",
		"--gateway",
		gateway,
		"--timeout",
		timeout.String(),
		"-f",
		yamlFile,
	}

	if tlsInsecure {
		args = append(args, "--tls-no-verify")
	}

	cmd := "faas-cli"

	return cmd, args
}
