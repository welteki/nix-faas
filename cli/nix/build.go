package nix

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func BuildStack(module string, gcRoot *GarbageCollectionRoot) error {
	var outLink string

	if gcRoot != nil {
		outLink = gcRoot.Path()
	}

	cmd, args := getBuildCommand(getEvalStackFile(), formatModulesExpr(module), outLink)

	task := execute.ExecTask{
		Command:     cmd,
		Args:        args,
		Shell:       true,
		StreamStdio: true,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("\"%s\" terminated with non-zero exit code: %s", cmd, res.Stderr)
	}

	return nil
}

func getBuildCommand(evalStackFile string, modulesExpr string, outLink string) (string, []string) {
	args := []string{
		evalStackFile,
		"--arg",
		"modules",
		modulesExpr,
		"--attr",
		"config.stackYaml",
	}

	if len(outLink) == 0 {
		args = append(args, "--no-out-link")
	} else {
		args = append(args, fmt.Sprintf("--out-link %s", outLink))
	}

	command := "nix-build"

	return command, args
}
