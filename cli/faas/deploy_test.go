package faas

import (
	"strings"
	"testing"
)

func Test_getDeployCommand(t *testing.T) {
	want := "deploy -f ./stack.yaml"

	wantCommand := "faas-cli"

	command, args := getDeployCommand("./stack.yaml")

	joined := strings.Join(args, " ")

	if joined != want {
		t.Errorf("getDeployCommand want: %q, got: %q", want, joined)
	}

	if command != wantCommand {
		t.Errorf("getDeployCommand want: %q, got: %q", wantCommand, command)
	}
}
