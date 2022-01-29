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
      version = versionNumber + "-" + versionSuffix;
      versionSuffix =
        if officialRelease
        then ""
        else "pre${builtins.substring 0 8 (self.lastModifiedDate)}-${self.shortRev or "dirty"}";
      commit = "${self.rev or "dirty"}";

      versionNumber = "0.1.0";
      officialRelease = false;

      supportedSystems = [ "x86_64-linux" ];

      overlays = [
        (final: prev: {
          nix-faas =
            let
              inherit (final) lib buildGoModule makeWrapper skopeo faas-cli;
            in
            buildGoModule {
              inherit version;
              pname = "nix-faas";

              src = ./.;
              subPackages = [ "cli" ];

              vendorSha256 = null;

              buildInputs = [ makeWrapper ];

              ldflags = [
                "-s"
                "-w"
                "-X github.com/welteki/nix-faas/cli/nix.NixDir=${placeholder "out"}/nix"
                "-X github.com/welteki/nix-faas/cli/commands.Version=${version}"
                "-X github.com/welteki/nix-faas/cli/commands.GitCommit=${commit}"
              ];

              postInstall = ''
                makeWrapper $out/bin/cli $out/bin/nix-faas \
                  --prefix PATH : ${lib.makeBinPath [ skopeo faas-cli ]}

                mkdir $out/nix
                cp -R lib pkgs modules $out/nix
                cp ./overlay.nix $out/nix/overlay.nix
              '';
            };
        })
        localPkgsOverlay
      ];
      localPkgsOverlay = import ./overlay.nix;
    in
    {
      overlay = localPkgsOverlay;

    } // utils.lib.eachSystem supportedSystems (system:
      let
        pkgs = import nixpkgs {
          inherit system overlays;
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
