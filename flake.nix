{
  description = "Build and deploy serverless functions with Nix";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-21.11-small";
    utils.url = "github:numtide/flake-utils";
    flake-compat.url = "github:edolstra/flake-compat";
    flake-compat.flake = false;
  };

  outputs = { self, nixpkgs, utils, ... }@inputs:
    let
      supportedSystems = [ "x86_64-linux" ];
    in
    {

      overlay = final: prev: {
        of-watchdog = import ./pkgs/of-watchdog.nix final;
        classic-watchdog = import ./pkgs/classic-watchdog.nix final;
        ofTools = import ./pkgs/build-support/openfaas final;

        nix-faas =
          let
            inherit (final) lib buildGoModule makeWrapper skopeo faas-cli;
          in
          buildGoModule {
            pname = "nix-faas";
            version = "0.0.1-dev";

            src = ./.;
            subPackages = [ "cli" ];

            vendorSha256 = null;

            buildInputs = [ makeWrapper ];

            ldflags = [
              "-s"
              "-w"
              "-X github.com/welteki/nix-faas/cli/nix.NixDir=${placeholder "out"}/nix"
            ];

            postInstall = ''
              makeWrapper $out/bin/cli $out/bin/nix-faas \
                --prefix PATH : ${lib.makeBinPath [ skopeo faas-cli ]}

              mkdir $out/nix
              cp -R lib pkgs modules $out/nix
            '';
          };
      };

    } // utils.lib.eachSystem supportedSystems (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ self.overlay ];
        };
      in
      {
        packages = {
          inherit (pkgs) nix-faas of-watchdog classic-watchdog;
          classic-watchdog-image = pkgs.ofTools.baseImages.classic-watchdog;
          of-watchdog-image = pkgs.ofTools.baseImages.of-watchdog;
        };

        defaultPackage = pkgs.nix-faas;

        devShell = pkgs.mkShell {
          buildInputs = builtins.attrValues {
            inherit (pkgs)
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
              nixpkgs-fmt;
          };
        };
      });
}
