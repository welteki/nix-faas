package nix

import (
	"os"
	"path"
	"strings"
	"testing"
)

func Test_getEvaluateCommand(t *testing.T) {
	want := "./lib/eval-stack.nix --eval --strict --json --arg modules \"[ ./stack.nix ]\" --attr config.stackYamlText"

	wantCommand := "nix-instantiate"

	command, args := getEvaluateCommand("./lib/eval-stack.nix", "\"[ ./stack.nix ]\"")

	joined := strings.Join(args, " ")

	if joined != want {
		t.Errorf("getEvaluateCommand want: \"%s\", got: \"%s\"", want, joined)
	}

	if command != wantCommand {
		t.Errorf("getEvaluateCommand want command: \"%s\", got: \"%s\"", wantCommand, command)
	}
}

func Test_getNixDir_EmptyNixDir(t *testing.T) {
	dir := "/nix/store/faas-nix"
	os.Setenv(NixDirEnv, dir)

	want := dir
	got := getNixDir()
	if want != got {
		t.Fatalf("getNixDir want %s env value: \"%s\", got: \"%s\"", NixDirEnv, want, got)
	}
}

func Test_getNixDir(t *testing.T) {
	dir := "/nix/store/faas-nix"
	NixDir = dir

	want := dir
	got := getNixDir()
	if want != got {
		t.Fatalf("getNixDir want: \"%s\", got: \"%s\"", want, got)
	}
}

func Test_getEvalStackFile(t *testing.T) {
	NixDir = "/nix/store/faas-nix"

	want := path.Join(NixDir, "lib/eval-stack.nix")
	got := getEvalStackFile()
	if want != got {
		t.Fatalf("getEvalStackFile want: \"%s\", got: \"%s\"", want, got)
	}
}

func Test_formatModuleExpr_Rel(t *testing.T) {
	want := "\"[ (./. + \"/stack.nix\") ]\""
	got := formatModulesExpr("stack.nix")
	if want != got {
		t.Fatalf("formatModuleExpr want: \"%s\", got: \"%s\"", want, got)
	}
}

func Test_formatModulesExpr_Abs(t *testing.T) {
	want := "\"[ (/. + \"/foo/stack.nix\") ]\""
	got := formatModulesExpr("/foo/stack.nix")
	if want != got {
		t.Fatalf("formatModuleExpr want: \"%s\", got: \"%s\"", want, got)
	}
}
