{ pkgs ? import <nixpkgs> { }
, modules ? [ ]
}:

let
  inherit (pkgs) lib buildGoModule;

  serviceStack = lib.evalModules {
    check = true;
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
  };
in
serviceStack
