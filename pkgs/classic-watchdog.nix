{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "classic-watchdog";
  version = "0.2.3";
  rev = "b20362a05bde813764e0b20d27de58df2159bc0a";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "classic-watchdog";
    rev = version;
    sha256 = "sha256-z3jvG/QUI815wyTo5IwcbGR3Pa56J2WNiTIZZ6FCGz0=";
  };

  vendorHash = null;

  CGO_ENABLED = 0;

  subPackages = [ "." ];

  ldflags = [
    "-s"
    "-w"
    "-X main.GitCommit=${rev}"
    "-X main.Version=${version}"
  ];

  postInstall = ''
    ln -s $out/bin/classic-watchdog $out/bin/fwatchdog
  '';
}
