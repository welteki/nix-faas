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
    x-nix-faas = config.of-stack.extended;
  };
  stackYamlText = builtins.toJSON stack;
  stackYaml = pkgs.writeText "stack.yaml" stackYamlText;
in
{
  options = {
    stackYaml = mkOption {
      description = "A derivation that produces a stack.yaml file for this function stack.";
      type = package;
      readOnly = true;
    };

    stackYamlText = mkOption {
      description = "Text of stack.yaml file.";
      type = str;
      readOnly = true;
    };

    provider.name = mkOption {
      description = "The only valid value for provider `name` is `openfaas`.";
      type = str;
      default = "openfaas";
    };

    provider.gateway = mkOption {
      description = "The gateway URL.";
      type = str;
      default = "http://127.0.0.1:8080";
    };

    of-stack.extended = mkOption {
      description = "Attribute set that will be turned into the x-nix-faas section of the stack.yaml file.";
      type = attrs;
    };

    functions = mkOption {
      description = "An attribute set of function configurations";
      type = attrsOf (submodule of-function);
      default = { };
      example = {
        figlet = {
          image = {
            name = "docker.io/welteki/figlet";
          };

          fprocess = ''''${pkgs.figlet}/bin/figlet'';
        };
      };
    };
  };

  config = {
    inherit stackYaml stackYamlText;
  };
}
