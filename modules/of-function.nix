{ pkgs, lib, config, ... }:

let
  inherit (lib) mkOption types;
  inherit (types) listOf nullOr attrsOf str either int bool enum package;
  inherit (pkgs) writeShellScriptBin dockerTools callPackage;

  classic-watchdog = callPackage ../pkgs/classic-watchdog.nix { };
  of-watchdog = callPackage ../pkgs/of-watchdog.nix { };

  defaultTo = default: v: if v == null then default else v;

  fwatchdog =
    if config.watchdog == "classic"
    then "${classic-watchdog}/bin/classic-watchdog"
    else "${of-watchdog}/bin/of-watchdog";

  function = writeShellScriptBin "${config.name}" ''
    export fprocess="${builtins.concatStringsSep " " config.fprocess}"
    export mode="${defaultTo "streaming" config.watchdogMode}"

    ${fwatchdog}
  '';

  functionImage = dockerTools.streamLayeredImage {
    inherit (config.image) name tag;

    extraCommands = ''
      mkdir -p tmp
      mkdir -p var/openfaas/secrets
    '';

    config.Cmd = [ "${function}/bin/${function.name}" ];
  };

  limitOptions = {
    cpu = mkOption {
      type = nullOr str;
      default = null;
    };

    memory = mkOption {
      type = nullOr str;
      default = null;
    };
  };
in

{
  options = {
    stackEntry = mkOption {
      type = attrsOf types.unspecified;
      readOnly = true;
    };

    build.image = mkOption {
      type = package;
      internal = true;
    };

    build.imageSpecifier = mkOption {
      type = str;
      internal = true;
    };

    name = mkOption {
      type = str;
      description = ''
        The name of the function - `<name>` in the stack-level `functions.<name>`
      '';
      readOnly = true;
    };

    fprocess = mkOption {
      type = listOf str;
      default = [ ];
    };

    watchdog = mkOption {
      type = enum [ "classic" "of" ];
      default = "classic";
    };

    watchdogMode = mkOption {
      type = enum [ "http" "streaming" "serializing" "static" ];
      default = "streaming";
    };

    image.name = mkOption {
      type = str;
      default = config.name;
    };

    image.tag = mkOption {
      type = nullOr str;
      default = null;
    };

    environment = mkOption {
      type = attrsOf (either str int);
      default = { };
    };

    secrets = mkOption {
      type = listOf str;
      default = [ ];
    };

    readonly_root_filesystem = mkOption {
      type = nullOr bool;
      default = null;
    };

    constraints = mkOption {
      type = listOf str;
      default = [ ];
    };

    labels = mkOption {
      type = attrsOf str;
      default = { };
    };

    annotations = mkOption {
      type = attrsOf str;
      default = { };
    };

    limits = limitOptions;

    requests = limitOptions;
  };

  config = {
    build.image = functionImage;
    build.imageSpecifier =
      "${config.image.name}:${
        if config.image.tag != null
        then config.image.tag
        else lib.head (lib.strings.splitString "-" (baseNameOf config.build.image.outPath))}";

    stackEntry = {
      image = config.build.imageSpecifier;

    } // lib.optionalAttrs (config.environment != { }) {
      inherit (config) environment;
    } // lib.optionalAttrs (config.secrets != [ ]) {
      inherit (config) secrets;
    } // lib.optionalAttrs (config.readonly_root_filesystem != null) {
      inherit (config) readonly_root_filesystem;
    } // lib.optionalAttrs (config.constraints != [ ]) {
      inherit (config) constraints;
    } // lib.optionalAttrs (config.labels != { }) {
      inherit (config) labels;
    } // lib.optionalAttrs (config.annotations != { }) {
      inherit (config) annotations;
    } // lib.optionalAttrs (config.limits.cpu != null || config.limits.memory != null) {
      inherit (config) limits;
    } // lib.optionalAttrs (config.requests.cpu != null || config.requests.memory != null) {
      inherit (config) requests;
    };
  };
}
