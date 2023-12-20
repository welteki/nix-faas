package commands

import (
	"fmt"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/google/go-containerregistry/pkg/v1/tarball"
	"github.com/spf13/cobra"

	"github.com/welteki/nix-faas/cli/image"
	"github.com/welteki/nix-faas/cli/nix"
	"github.com/welteki/nix-faas/cli/stack"
)

func init() {
	rootCmd.AddCommand(publishCmd)
}

var publishCmd = &cobra.Command{
	Use:     "publish -f MODULE_FILE",
	Short:   "Builds and pushes OpenFaaS function images to a container registry",
	PreRunE: preRunPublish,
	RunE:    runPublish,
}

func preRunPublish(cmd *cobra.Command, args []string) error {
	if len(stackModules) == 0 {
		return fmt.Errorf("use --file or -f to specify modules")
	}

	return nil
}

func runPublish(cmd *cobra.Command, args []string) (retErr error) {
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

	return nil
}

func push(m stack.ImageMetadata) (retErr error) {
	dockerArchive, err := image.NewArchiveFromStream(m.Source)
	if err != nil {
		return fmt.Errorf("creating image archive: %w", err)
	}
	defer func() {
		if err := dockerArchive.Close(); err != nil {
			retErr = fmt.Errorf("(archive: %v): %w", err, retErr)
		}
	}()

	ref, err := name.ParseReference(m.Specifier)
	if err != nil {
		return err
	}

	img, err := tarball.ImageFromPath(dockerArchive.Path(), nil)
	if err != nil {
		return err
	}

	options := []remote.Option{
		remote.WithAuthFromKeychain(authn.DefaultKeychain),
	}
	if err := remote.Write(ref, img, options...); err != nil {
		return err
	}

	return nil
}
