package faas

import (
	"strings"
	"testing"
	"time"
)

func Test_getDeployCommand(t *testing.T) {
	want := "deploy --gateway https://gw.example.com --timeout 1m0s -f ./stack.yaml"

	wantCommand := "faas-cli"

	command, args := getDeployCommand("./stack.yaml", "https://gw.example.com", 60*time.Second, false)

	joined := strings.Join(args, " ")

	if joined != want {
		t.Errorf("getDeployCommand want: %q, got: %q", want, joined)
	}

	if command != wantCommand {
		t.Errorf("getDeployCommand want: %q, got: %q", wantCommand, command)
	}
}

func Test_getDeployCommand_tlsInsecure(t *testing.T) {
	want := "deploy --gateway https://gw.example.com --timeout 1m0s -f ./stack.yaml --tls-no-verify"

	wantCommand := "faas-cli"

	command, args := getDeployCommand("./stack.yaml", "https://gw.example.com", 60*time.Second, true)

	joined := strings.Join(args, " ")

	if joined != want {
		t.Errorf("getDeployCommand want: %q, got: %q", want, joined)
	}

	if command != wantCommand {
		t.Errorf("getDeployCommand want: %q, got: %q", wantCommand, command)
	}
}
