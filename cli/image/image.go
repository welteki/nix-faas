package image

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

// Archive represents a (docker-formatted) tar archive of an image.
type Archive struct {
	path          string
	removeOnClose bool
}

// NewArchiveFromFile returns an Archive for the specified path.
// The caller should call .Close() on the returned archive when done.
func NewArchiveFromFile(path string) (*Archive, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer file.Close()

	return newArchive(file.Name(), false)
}

// NewArchiveFromStream returns an Archive for the specified (nix streamLayeredImage generated) image build script.
// The caller should call .Close() on the returned archive when done.
func NewArchiveFromStream(script string) (*Archive, error) {
	tarArchive, err := ioutil.TempFile("", "nix-faas-docker-tar-*")
	if err != nil {
		return nil, fmt.Errorf("creating temporary file: %w", err)
	}
	defer tarArchive.Close()

	succeeded := false
	defer func() {
		if !succeeded {
			os.Remove(tarArchive.Name())
		}
	}()

	if err := generateImage(script, tarArchive); err != nil {
		return nil, fmt.Errorf("generating image: %w", err)
	}
	succeeded = true

	return newArchive(tarArchive.Name(), true)
}

// newArchive creates an Archive for the specified path and removeOnClose flag.
// The caller should call .Close() on the returned archive when done.
func newArchive(path string, removeOnClose bool) (*Archive, error) {
	i := Archive{
		path:          path,
		removeOnClose: removeOnClose,
	}

	return &i, nil
}

// Path returns the path to the image archive.
func (i *Archive) Path() string {
	return i.path
}

// Close removes resources associated with an initialized Archive, if any.
func (i *Archive) Close() error {
	path := i.path
	i.path = ""
	if i.removeOnClose {
		return os.Remove(path)
	}
	return nil
}

// generateImage generates an image archive by executing the specified (nix streamLayeredImage generated)
// image build script and streaming the output to the specified writer.
func generateImage(e string, o io.Writer) error {
	cmd := exec.Command(e)
	cmd.Stdout = o

	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("starting generate script: %w", err)
	}

	err = cmd.Wait()
	if err != nil {
		return fmt.Errorf("executing generate script %s: %w", e, err)
	}

	return nil
}
