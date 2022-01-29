{ pkgs, lib, config, ... }:

let
  inherit (lib) mkOption types;
  inherit (types) listOf nullOr attrsOf str either int bool enum package;

  functionImage = pkgs.ofTools.streamOfImage {
    inherit (config.image) name tag created;
    inherit (config) watchdog watchdogMode fprocess;
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
      type = nullOr str;
      default = null;
    };

    watchdog = mkOption {
      type = enum [ "classic-watchdog" "of-watchdog" ];
      default = "classic-watchdog";
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

    image.created = mkOption {
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
