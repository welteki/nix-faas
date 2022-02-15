{ pkgs, lib, config, ... }:

let
  inherit (lib) mkOption types;
  inherit (types) listOf nullOr attrsOf str either int bool enum package;

  ref = url: text:
    ''See: [${text}](${url})'';

  ofDocsRef = path: text:
    ref ''https://docs.openfaas.com/${path}'' text;
  ofYamlRef = fragment:
    ofDocsRef ''/reference/yaml/#${fragment}'' "OpenFaaS YAML file reference";

  functionImage = pkgs.ofTools.streamOfImage {
    inherit (config.image) name tag created;
    inherit (config) watchdog watchdogMode fprocess;
  };

  limitOptions = {
    cpu = mkOption {
      description = ''
        CPU limit/request for a function.

        ${ofYamlRef "function-memorycpu-limits"}
      '';
      type = nullOr str;
      default = null;
      example = "100m";
    };

    memory = mkOption {
      description = ''
         Memory limit/request for a function.

        ${ofYamlRef "function-memorycpu-limits"}
      '';
      type = nullOr str;
      default = null;
      example = "40Mi";
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
      description = ''
        Process to execute.

        When using the `of-watchdog` in `http` mode this process is executed to start a server.
        For non http modes or when using the `classic-watchdog` this process is executed on each request. In this case
        it must accept inputs via STDIN and print outputs via STDOUT.

        ${ofDocsRef "architecture/watchdog/" "OpenFaas watchdog architecture"}
      '';
      type = nullOr str;
      default = null;
      example = ''''${pkgs.figlet}/bin/figlet'';
    };

    watchdog = mkOption {
      description = ''
        Watchdog used for the function.

        ${ofDocsRef "architecture/watchdog/" "OpenFaas watchdog architecture"}
      '';
      type = enum [ "classic-watchdog" "of-watchdog" ];
      default = "classic-watchdog";
    };

    watchdogMode = mkOption {
      description = ''
        Watchdog mode used for the function. (Only applicable when using the `of-watchdog`).

        ${ref "https://github.com/openfaas/of-watchdog/blob/master/README.md" "of-watchdog README"}
      '';
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
      description = ''
        Environment variable to set in the function.
      '';
      type = attrsOf (either str int);
      default = { };
      example = {
        http_proxy = "http://proxy1.corp.com:3128";
        no_proxy = "http://gateway/";
      };
    };

    secrets = mkOption {
      description = ''
        List of secrets to make available in the function.

        OpenFaaS functions can make use of secure secrets using the secret store from Kubernetes or faasd.
        All secrets are made available in the container file-system and should be read from `/var/openfaas/secrets/<secret-name>`.

        ${ofDocsRef "reference/secrets/" "OpenFaaS secrets reference"}
      '';
      type = listOf str;
      default = [ ];
      example = [
        "s3_access_key"
        "s3_secret_key"
      ];
    };

    readonly_root_filesystem = mkOption {
      description = ''
        The readonly_root_filesystem indicates that the function file system will be set to read-only except for the temporary folder /tmp.

        See: https://docs.openfaas.com/reference/yaml/#function-read-only-root-filesystem
      '';
      type = bool;
      default = false;
    };

    constraints = mkOption {
      description = ''
        Set contraints on functions. eg: assign functions to a given NodePool or Node, pin functions to a certain host or type of host.

        ${ofYamlRef "function-constraints"}
      '';
      type = listOf str;
      default = [ ];
      example = [
        "node.platform.os == linux"
      ];
    };

    labels = mkOption {
      description = ''
        Add labeles to the function.

        Labels are passed directly to the container scheduler. They are also available from the OpenFaaS REST API for quering or grouping functions.

        ${ofYamlRef "function-labels"}
      '';
      type = attrsOf str;
      default = { };
      example = {
        canary = true;
        Git-Owner = "alexellis";
      };
    };

    annotations = mkOption {
      description = ''
        Add annotations to the function.

        Annotations are a collection of meta-data which is stored with the function by the provider. Annotations are also available from the OpenFaaS REST API for querying.

        ${ofYamlRef "function-annotations"}
      '';
      type = attrsOf str;
      default = { };
      example = {
        "com.openfaas.health.http.path" = "/healtz";
        "com.openfaas.health.http.initialDelay" = "30s";
      };
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
