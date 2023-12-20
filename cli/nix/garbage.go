package nix

import (
	"fmt"
	"os"
	"path/filepath"
)

// GarbageCollectionRoot represents a nix garbage collection root.
// The path of the garbage collection root will be passed to the "--add-root" option during the
// the realisition of nix store paths. Calling .Close() on the garbage collection root will remove its associated
// resources and thus allow the store paths to be deleted by the nix garbage collector.
// See: https://nixos.org/manual/nix/stable/command-ref/nix-store.html#common-options
type GarbageCollectionRoot struct {
	path          string
	gcPath        string
	removeOnClose bool
}

// NewGarbageCollectionRoot returns a GarbageCollectionRoot.
// The caller should call .Close() on the returned garbage collection root when done.
func NewGarbageCollectionRoot() (*GarbageCollectionRoot, error) {
	gcDir, err := os.MkdirTemp("", "nix-faas-gcroot-*")
	if err != nil {
		return nil, err
	}

	return &GarbageCollectionRoot{
		path:          filepath.Join(gcDir, "result"),
		gcPath:        gcDir,
		removeOnClose: true,
	}, nil
}

// NewGarbageCollectionRootFromPath returns a GarbageCollectionRoot for the specified path.
// The caller should call .Close() on the returned garbage collection root when done.
// When preserve is true, resources associated with the garbage collection root will be preserved even
// when calling .Close()
func NewGarbageCollectionRootFromPath(path string, preserve bool) (*GarbageCollectionRoot, error) {
	_, err := os.Stat(path)

	if err == nil {
		return nil, fmt.Errorf("invalid garbage collection root %q: path already exists", path)
	}

	return &GarbageCollectionRoot{
		path:          path,
		gcPath:        path,
		removeOnClose: preserve,
	}, nil
}

// Close removes resources associated with an initialized GarbageCollectionRoot, if any.
func (gc *GarbageCollectionRoot) Close() error {
	gc.path = ""
	if gc.removeOnClose {
		return os.RemoveAll(gc.gcPath)
	}
	return nil
}

// Path returns the path of the garbage collection root.
func (gc *GarbageCollectionRoot) Path() string {
	return gc.path
}
