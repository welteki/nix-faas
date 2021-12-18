module github.com/welteki/nix-faas

go 1.16

require (
	github.com/alexellis/go-execute v0.5.0
	github.com/spf13/cobra v1.2.1
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
	gopkg.in/yaml.v2 v2.4.0
)

// Setting Shell option for execute does not work since it assumes /bin/bash is present on
// the system which is not the case for NixOS.
// Replace go-execute with patched version.
replace github.com/alexellis/go-execute => github.com/welteki/go-execute v0.5.1-0.20211218163025-d60c089639ba
