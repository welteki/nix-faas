{ pkgs ? import <nixpkgs> { overlays = [ (import ../overlay.nix) ]; }
, modules ? [ ]
}:

let
  inherit (pkgs) lib;

  serviceStack = lib.evalModules {
    modules = baseModules ++ modules;
  };

  baseModules = [
    pkgsModule 
  ] ++ import ../modules;

  pkgsModule = {
    _file = ./eval-stack.nix;
    key = ./eval-stack.nix;

    config._module.args.pkgs = lib.mkIf (pkgs != null) (lib.mkForce pkgs);
    config._module.check = true;
  };
in
serviceStack
