package commands

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	execute "github.com/alexellis/go-execute/pkg/v1"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(evalCmd)
}

var evalCmd = &cobra.Command{
	Use:     "eval -f MODULE_FILE",
	Short:   "Show the stack file for the current configuration as JSON",
	PreRunE: preRunEval,
	RunE:    runEval,
}

func preRunEval(cmd *cobra.Command, args []string) error {
	if len(stackModule) == 0 {
		return fmt.Errorf("--file or -f is required")
	}

	return nil
}

func runEval(cmd *cobra.Command, args []string) error {
	stackYaml, err := eval(stackModule)

	if err != nil {
		return err
	}

	fmt.Printf("%s", stackYaml)
	return nil
}

func eval(module string) (string, error) {
	cmd := "nix-instantiate"

	args := []string{
		path.Join(getNixDir(), "lib/eval-stack.nix"),
		"--eval",
		"--read-write-mode",
		"--show-trace",
		"--json",
		fmt.Sprintf("--arg modules \"[ \"%s\" ]\"", module),
		"--attr config.stackYamlText",
	}

	// Setting Shell option for execute does not work since it assumes /bin/bash is present on
	// the system which is not the case for NixOS.
	// TODO: raise issue upstream
	script := strings.Join(args, " ")

	task := execute.ExecTask{
		Command: "bash",
		Args:    append([]string{"-c"}, fmt.Sprintf("%s %s", cmd, script)),
	}

	res, err := task.Execute()

	if err != nil {
		return "", err
	}

	if res.ExitCode != 0 {
		return "", fmt.Errorf("received not-zero exit code from evaluation, error: %s", res.Stderr)
	}

	stackYaml, _ := strconv.Unquote(res.Stdout)

	return stackYaml, nil
}
