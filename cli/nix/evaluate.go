package nix

import (
	"fmt"
	"os"
	"path"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

const NixDirEnv string = "NIXFAAS_NIX_DIR"

// Path to nix files
var NixDir string

func EvaluateStack(module string) (string, error) {
	cmd, args := getEvaluateCommand(getEvalStackFile(), formatModulesExpr(module))

	task := execute.ExecTask{
		Command: cmd,
		Args:    args,
		Shell:   true,
	}

	res, err := task.Execute()

	if err != nil {
		return "", err
	}

	if res.ExitCode != 0 {
		return "", fmt.Errorf("\"%s\" terminated with non-zero exit code: %s", cmd, res.Stderr)
	}

	return res.Stdout, nil
}

func getEvaluateCommand(evalStackFile string, modulesExpr string) (string, []string) {
	args := []string{
		evalStackFile,
		"--eval",
		"--strict",
		"--json",
		"--arg",
		"modules",
		modulesExpr,
		"--attr",
		"config.stackYamlText",
	}

	command := "nix-instantiate"

	return command, args
}

func getNixDir() string {
	if len(NixDir) == 0 {
		return os.Getenv(NixDirEnv)
	}

	return NixDir
}

func getEvalStackFile() string {
	return path.Join(getNixDir(), "lib/eval-stack.nix")
}

func formatModulesExpr(module string) string {
	if path.IsAbs(module) {
		return fmt.Sprintf("\"[ (/. + \"%s\") ]\"", module)
	} else {
		return fmt.Sprintf("\"[ (./. + \"/%s\") ]\"", module)
	}
}
