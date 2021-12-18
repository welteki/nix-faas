package commands

import (
	"fmt"
	"path"

	"github.com/spf13/cobra"

	execute "github.com/alexellis/go-execute/pkg/v1"
)

var (
	outLink string
)

func init() {
	buildCmd.Flags().StringVarP(&outLink, "out-link", "o", "", "Create a symlink from the derivations output path to outlink")

	rootCmd.AddCommand(buildCmd)
}

var buildCmd = &cobra.Command{
	Use:     "build -f MODULE_FILE",
	Short:   "Build OpenFaaS functions",
	PreRunE: preRunBuild,
	RunE:    runBuild,
}

func preRunBuild(cmd *cobra.Command, args []string) error {
	if len(stackModule) == 0 {
		return fmt.Errorf("--file or -f is required")
	}

	return nil
}

func runBuild(cmd *cobra.Command, args []string) error {
	err := build(stackModule, outLink, false)

	if err != nil {
		return err
	}

	return nil
}

func build(module string, outLink string, quitBuild bool) error {
	cmd := "nix-build"

	args := []string{
		path.Join(getNixDir(), "lib/eval-stack.nix"),
		fmt.Sprintf("--arg modules \"[ \"%s\" ]\"", module),
		"--show-trace",
		"--attr config.stackYaml",
	}

	if len(outLink) == 0 {
		args = append(args, "--no-out-link")
	} else {
		args = append(args, fmt.Sprintf("--out-link %s", outLink))
	}

	task := execute.ExecTask{
		Command:     cmd,
		Args:        args,
		Shell:       true,
		StreamStdio: !quitBuild,
	}

	res, err := task.Execute()

	if err != nil {
		return err
	}

	if res.ExitCode != 0 {
		return fmt.Errorf("received non-zero exit code from build, error: %s", res.Stderr)
	}

	return nil
}
