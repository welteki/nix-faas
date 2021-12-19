package nix

import (
	"strings"
	"testing"
)

func Test_getBuildCommand_NoOutLink(t *testing.T) {
	want := "./lib/eval-stack.nix --arg modules \"[ ./stack.nix ]\" --attr config.stackYaml --no-out-link"

	wantCommand := "nix-build"

	command, args := getBuildCommand("./lib/eval-stack.nix", "\"[ ./stack.nix ]\"", "")

	joined := strings.Join(args, " ")

	if joined != want {
		t.Errorf("getBuildCommand want: \"%s\", got: \"%s\"", want, joined)
	}

	if command != wantCommand {
		t.Errorf("getBuildCommand want command: \"%s\", got: \"%s\"", wantCommand, command)
	}
}

func Test_getBuildCommand_OutLink(t *testing.T) {
	want := "./lib/eval-stack.nix --arg modules \"[ ./stack.nix ]\" --attr config.stackYaml --out-link out"

	wantCommand := "nix-build"

	command, args := getBuildCommand("./lib/eval-stack.nix", "\"[ ./stack.nix ]\"", "out")

	joined := strings.Join(args, " ")

	if joined != want {
		t.Errorf("getBuildCommand want: \"%s\", got: \"%s\"", want, joined)
	}

	if command != wantCommand {
		t.Errorf("getBuildCommand want command: \"%s\", got: \"%s\"", wantCommand, command)
	}
}
