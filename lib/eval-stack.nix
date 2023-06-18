{ pkgs ? import <nixpkgs> { overlays = [ (import ../overlay.nix) ]; }
, modules ? [ ]
}:

let
  inherit (pkgs) lib buildGoModule;

  serviceStack = lib.evalModules {
    modules = baseModules ++ modules;
  };

  baseModules = [ pkgsModule ] ++ [
    ../modules/of-stack.nix
    ../modules/images.nix
  ];

  pkgsModule = rec {
    _file = ./eval-stack.nix;
    key = _file;

    config._module.args.pkgs = lib.mkIf (pkgs != null) (lib.mkForce pkgs);
    config._module.check = true;
  };
in
serviceStack
