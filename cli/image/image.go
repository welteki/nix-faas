package image

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
)

type Archive struct {
	path          string
	removeOnClose bool
}

func NewArchiveFromFile(path string) (*Archive, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("opening file %q: %w", path, err)
	}
	defer file.Close()

	return newImageArchive(file.Name(), false)
}

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

	return newImageArchive(tarArchive.Name(), true)
}

func newImageArchive(path string, removeOnClose bool) (*Archive, error) {
	i := Archive{
		path:          path,
		removeOnClose: removeOnClose,
	}

	return &i, nil
}

func (i *Archive) Path() string {
	return i.path
}

func (i *Archive) Close() error {
	path := i.path
	i.path = ""
	if i.removeOnClose {
		return os.Remove(path)
	}
	return nil
}

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
