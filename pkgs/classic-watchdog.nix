{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "classic-watchdog";
  version = "0.2.2";
  rev = "f997bb97fea3a3196eb054994f2d44b038b297c2";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "classic-watchdog";
    rev = version;
    sha256 = "sha256-73cEGeP0trRAOeEDayTXG5mK25K3i5L1N1m02h+QXuA=";
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
