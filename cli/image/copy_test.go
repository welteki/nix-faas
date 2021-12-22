package image

import (
	"strings"
	"testing"
)

func Test_getCopyCommand(t *testing.T) {
	want := "copy docker-archive:archive/path docker://registry/nix-faas:latest --insecure-policy"

	wantCommand := "skopeo"

	command, args := getCopyCommand("archive/path", "registry/nix-faas:latest")

	joined := strings.Join(args, " ")

	if joined != want {
		t.Errorf("getCopyCommand want: %q, got: %q", want, joined)
	}

	if command != wantCommand {
		t.Errorf("getCopyCommand want: %q, got: %q", wantCommand, command)
	}
}
