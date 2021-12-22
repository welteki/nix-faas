package image

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

func Copy(srcRef, destRef string) error {
	cmd, args := getCopyCommand(srcRef, destRef)

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
		return fmt.Errorf("%q terminated with non-zero exit code: %s", cmd, res.Stderr)
	}

	return nil
}

func getCopyCommand(srcRef, destRef string) (string, []string) {
	args := []string{
		"copy",
		fmt.Sprintf("docker-archive:%s", srcRef),
		fmt.Sprintf("docker://%s", destRef),
		"--insecure-policy",
	}

	cmd := "skopeo"

	return cmd, args
}
