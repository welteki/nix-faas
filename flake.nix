{
  description = "Build and deploy serverless functions with Nix";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-21.11-small";
    utils.url = "github:numtide/flake-utils";
    flake-compat.url = "github:edolstra/flake-compat";
    flake-compat.flake = false;
  };

  outputs = { self, nixpkgs, utils, ... }@inputs: {

    overlay = final: prev: {
      of-watchdog = import ./pkgs/of-watchdog.nix final;
      classic-watchdog = import ./pkgs/classic-watchdog.nix final;
      ofTools = import ./pkgs/build-support/openfaas final;
    };

  } // utils.lib.eachSystem [ "x86_64-linux" ] (system:
    let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ self.overlay ];
      };
    in
    {
      packages = {
        inherit (pkgs) of-watchdog classic-watchdog;
        classic-watchdog-image = pkgs.ofTools.baseImages.classic-watchdog;
        of-watchdog-image = pkgs.ofTools.baseImages.of-watchdog;
      };

      devShell = pkgs.mkShell {
        buildInputs = [ pkgs.nixpkgs-fmt ];
      };
    });
}
