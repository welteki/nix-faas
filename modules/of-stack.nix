evalArguments@{ pkgs, lib, config, ... }:
let
  inherit (lib) mkOption mapAttrs;
  inherit (lib.types) attrsOf attrs str submodule package;

  of-function = {
    imports = [ argsModule ./of-function.nix ];
  };

  argsModule =
    { name # injected by types.submodule
    , ...
    }: rec {
      _file = ./of-stack.nix;
      key = _file;

      config._module.args.pkgs = lib.mkDefault evalArguments.pkgs;
      config.name = name;
    };

  stack = {
    inherit (config) provider;
    functions = mapAttrs (k: c: c.stackEntry) config.functions;
    nix-faas = config.of-stack.extended;
  };

  stackYaml = pkgs.writeText "stack.yaml" (builtins.toJSON stack);
in
{
  options = {
    stackYaml = mkOption {
      type = package;
      description = "A derivation that produces a stack.yaml file for this function stack.";
      readOnly = true;
    };

    provider.name = mkOption {
      type = str;
      default = "openfaas";
    };

    provider.gateway = mkOption {
      type = str;
      default = "http://127.0.0.1:8080";
    };

    of-stack.extended = mkOption {
      type = attrs;
    };

    functions = mkOption {
      type = attrsOf (submodule of-function);
      description = "An attribute set of function configurations";
    };
  };

  config = {
    inherit stackYaml;
  };
}
