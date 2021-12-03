{
  description = "Build and deploy serverless functions with Nix";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-23.11-small";
    utils.url = "github:numtide/flake-utils";
    flake-compat.url = "github:edolstra/flake-compat";
    flake-compat.flake = false;
  };

  nixConfig = {
    extra-substituters = [ "https://welteki.cachix.org" ];
    extra-trusted-public-keys = [ "welteki.cachix.org-1:zb0txiNEbjq9Fx7svp4LhTgFIQHKSa5ESi7QlLFjjQY=" ];
  };

  outputs = { self, nixpkgs, utils, ... }@inputs: {

    overlays.default = final: prev: {
      of-watchdog = import ./pkgs/of-watchdog.nix final;
      classic-watchdog = import ./pkgs/classic-watchdog.nix final;
      ofTools = import ./pkgs/build-support/openfaas final;

      hello-fn = import ./examples/hello-fn.nix final;
    };

  } // utils.lib.eachSystem [ "x86_64-linux" ] (system:
    let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ self.overlays.default ];
      };
    in
    {
      packages = {
        inherit (pkgs) of-watchdog classic-watchdog hello-fn;
        classic-watchdog-image = pkgs.ofTools.baseImages.classic-watchdog;
        of-watchdog-image = pkgs.ofTools.baseImages.of-watchdog;
      };

      devShells.default = pkgs.mkShell {
        buildInputs = with pkgs; [
          # Golang
          go
          gotools
          gopls
          go-outline
          gocode
          gopkgs
          gocode-gomod
          godef
          golint
          delve
          
          # Nix
          cachix
          statix
          vulnix
          deadnix
          nil
        ];
      };
    });
}
