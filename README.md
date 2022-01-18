> 🛠 **Status: Early experimental phase**
>
> This project is very much a work in progress. What you see here is only an early preview of the project.

## Part one: nix openfaas tools
Build [OpenFaas](https://www.openfaas.com/) function images with Nix.
```nix
{pkgs, ...}:
let
    inherit (pkgs.pkgsStatic) figlet;
in
pkgs.ofTools.buildOfImage {
    name = "figlet";
    tag = "latest";

    watchdog = "of-watchdog";
    watchdogMode = "serializing";

    fprocess = "${figlet}/bin/figlet";
}
```
This example recreates the figlet image from the openfaas store but uses the `of-watchdog` instead of the `classic-watchdog` used in the original image.

## Part two: CLI and NixOS modules
> Not yet part of the repo but hope to publish it soon.

Instead of configuring functions in YAML files `nix-faas` will use the Nix language to build and configure functions (similar to what [Arion](https://github.com/hercules-ci/arion) does for docker-compose).

The `nix-faas` cli will be a wrapper around the [faas-cli](https://github.com/openfaas/faas-cli) making it possible to deploy Nix defined function to any system running [OpenFaas](https://www.openfaas.com/).
