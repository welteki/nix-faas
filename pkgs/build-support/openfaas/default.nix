{ pkgs, lib, classic-watchdog, of-watchdog, dockerTools, ... }:

with builtins;

let
  baseImages =
    let
      bash = pkgs.pkgsStatic.bash;

      extraCommands = ''
        mkdir -p tmp
        mkdir -p var/openfaas/secrets
      '';
    in
    {
      classic-watchdog = dockerTools.buildImage {
        inherit extraCommands;

        name = classic-watchdog.pname;
        tag = classic-watchdog.version;

        copyToRoot = pkgs.buildEnv {
          name = "image-root";
          paths = [ classic-watchdog bash ];
          pathsToLink = [ "/bin" ];
        };
      };

      of-watchdog = dockerTools.buildImage {
        inherit extraCommands;

        name = of-watchdog.pname;
        tag = of-watchdog.version;

        copyToRoot = pkgs.buildEnv {
          name = "image-root";
          paths = [ of-watchdog bash ];
          pathsToLink = [ "/bin" ];
        };
      };
    };

  buildOrStreamImage =
    { stream
    , watchdog ? "classic-watchdog"
    , watchdogMode ? "streaming"
    , fprocess ? null
    , config ? { }
    , ...
    }@args:

      assert (!elem watchdog [ "classic-watchdog" "of-watchdog" ]) ->
        trace "unsupported watchdog `${watchdog}`, must be one of `classic-watchdog` or `of-watchdog`" false;
      assert (!elem watchdogMode [ "http" "streaming" "serializing" "static" ]) ->
        trace "unknown watchdogMode `${watchdogMode}`, must be one of `http`, `serializing`, `streaming` or `static`" false;

      let
        fwatchdog = baseImages.${watchdog};

        buildImage = if stream then dockerTools.streamLayeredImage else dockerTools.buildLayeredImage;
        acceptedArgs = functionArgs dockerTools.streamLayeredImage;
        args' = intersectAttrs acceptedArgs args;
      in
      buildImage (args' // {
        fromImage = fwatchdog;

        # TODO: better way for merging config?
        config = config // {
          Cmd = [ "/bin/fwatchdog" ];

          Env = lib.optionals (config ? Env) config.Env
            ++ [ "mode=${watchdogMode}" ]
            ++ lib.optionals (fprocess != null) [ "fprocess=${fprocess}" ];

          ExposedPorts = {
            "8080" = { };
          };

          Healthcheck = {
            Test = [
              "CMD"
              "/bin/sh"
              "-c"
              "[ -e /tmp/.lock ] || exit 0"
            ];
            Interval = 3000000000;
          };
        };
      });
in
{
  inherit baseImages;

  buildOfImage = args: buildOrStreamImage (args // { stream = false; });
  streamOfImage = args: buildOrStreamImage (args // { stream = true; });
}
