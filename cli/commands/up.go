package commands

import (
	"fmt"
	"time"

	"github.com/spf13/cobra"
	"github.com/welteki/nix-faas/cli/faas"
	"github.com/welteki/nix-faas/cli/nix"
	"github.com/welteki/nix-faas/cli/stack"
)

const (
	defaultGateway     = "http://127.0.0.1:8080"
	faasTimeoutDefault = 60 * time.Second
)

var (
	gateway     string
	tlsInsecure bool
	faasTimeout time.Duration
)

func init() {
	upCmd.Flags().StringVarP(&gateway, "gateway", "g", defaultGateway, "Gateway URL starting with http(s)://")
	upCmd.Flags().DurationVar(&faasTimeout, "timeout", faasTimeoutDefault, "Timeout for HTTP calls to the OpenFaaS API.")
	upCmd.Flags().BoolVar(&tlsInsecure, "tls-no-verify", false, "Disable TLS validation for HTTP calls to the OpenFaaS API")

	rootCmd.AddCommand(upCmd)
}

var upCmd = &cobra.Command{
	Use:     "up -f MODULE_FILE",
	Short:   "Build, push and deploy OpenFaaS functions",
	PreRunE: preRunUp,
	RunE:    runUp,
}

func preRunUp(cmd *cobra.Command, args []string) error {
	if len(stackModules) == 0 {
		return fmt.Errorf("use --file or -f to specify modules")
	}

	return nil
}

func runUp(cmd *cobra.Command, args []string) (retErr error) {
	gcRoot, err := nix.NewGarbageCollectionRoot()
	if err != nil {
		return err
	}
	defer func() {
		if err := gcRoot.Close(); err != nil {
			retErr = fmt.Errorf("(gcroot: %v): %w", err, retErr)
		}
	}()

	err = nix.BuildStack(stackModules, gcRoot)
	if err != nil {
		return err
	}

	stackYaml := gcRoot.Path()

	config, err := stack.ReadNixFaasConfig(stackYaml)
	if err != nil {
		return fmt.Errorf("getting nix-faas config: %w", err)
	}

	for _, image := range config.StackMetadata.Images {
		err := push(image)
		if err != nil {
			return err
		}
	}

	err = faas.Deploy(stackYaml, gateway, faasTimeout, tlsInsecure)
	if err != nil {
		return fmt.Errorf("deploying functions: %w", err)
	}

	return nil
}
