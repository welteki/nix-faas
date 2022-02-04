package commands

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
	"github.com/welteki/nix-faas/cli/nix"
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
	if len(stackModules) == 0 {
		return fmt.Errorf("use --file or -f to specify modules")
	}

	return nil
}

func runEval(cmd *cobra.Command, args []string) error {
	stackYaml, err := nix.EvaluateStack(stackModules)

	if err != nil {
		return err
	}

	stackYaml, _ = strconv.Unquote(stackYaml)

	fmt.Printf("%s", stackYaml)
	return nil
}
