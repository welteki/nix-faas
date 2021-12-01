{
  description = "Flake with some boilerplate";

  inputs = {
    nixpkgs.follows = "nix/nixpkgs";
    utils.url = "github:numtide/flake-utils";
    flake-compat.url = "github:edolstra/flake-compat";
    flake-compat.flake = false;
  };

  outputs = { nix, self, nixpkgs, utils, ... }@inputs: {

    overlay = final: prev: {
      of-watchdog = import .packages/of-watchdog.nix final;
      classic-watchdog = import ./packages/classic-watchdog.nix final;
    };

  } // utils.lib.eachDefaultSystem (system:
    let
      pkgs = import nixpkgs {
        inherit system;
        overlays = [ self.overlay ];
      };
    in
    {
      packages = {
        inherit (pkgs) of-watchdog classic-watchdog;
      };

      devShell = pkgs.mkShell {
        buildInputs = [ pkgs.nixpkgs-fmt ];
      };
    });
}
