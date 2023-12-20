> 🛠 **Status: Early experimental phase**
>
> This project is very much a work in progress.

## Part one: Nix OpenFaaS tools

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

This example recreates the figlet image from the OpenFaaS store but uses the `of-watchdog` instead of the `classic-watchdog` used in the original image.

## Part two: CLI and NixOS modules

`nix-faas` is a tool to build and deploy [OpenFaaS](https://www.openfaas.com/) functions using NixOS modules.

Instead of configuring functions in YAML files `nix-faas` will use the Nix language to build, configure and deploy OpenFaaS functions (similar to what [Arion](https://github.com/hercules-ci/arion) does for docker-compose).

### Preview

Deploy the `hello` and `figlet` package as a function to OpenFaaS.

Example `stack.nix` file:

```nix
{pkgs, ...}: {
  functions = {
    hello = {
        fprocess = "${pkgs.pkgsStatic.hello}/bin/hello";
        image = {
            name = "ttl.sh/nix-faas/hello";
        };
    };

    figlet = {
        fprocess = "${pkgs.pkgsStatic.figlet}/bin/figlet";
        image = {
            name = "ttl.sh/nix-faas/figlet";
        };
    };
  };
}
```

Deploy:

```bash
nix-faas up -f stack.nix
```

![](docs/preview.gif)
