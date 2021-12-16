{ pkgs, lib, config, ... }:
let
  inherit (lib) mkOption types mapAttrs attrValues;
  inherit (types) listOf unspecified;

  functionImages = mapAttrs addDetails config.functions;

  addDetails = functionName: function:
    let
      inherit (function) build;
    in
    {
      specifier = build.imageSpecifier;
      source = build.image.outPath;
    };
in
{
  config = {
    of-stack.extended.images = attrValues functionImages;
  };
}
