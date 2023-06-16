{ pkgs, ofTools, ... }:

ofTools.buildOfImage {
  name = "hello-fn";

  watchdog = "of-watchdog";
  watchdogMode = "serializing";

  fprocess = "${pkgs.pkgsStatic.hello}/bin/hello";
}
