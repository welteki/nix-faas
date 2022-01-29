final: prev: {
  of-watchdog = import ./pkgs/of-watchdog.nix final;
  classic-watchdog = import ./pkgs/classic-watchdog.nix final;
  ofTools = import ./pkgs/build-support/openfaas final;

  hello-fn = import ./examples/hello-fn.nix final;
}
