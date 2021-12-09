{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "classic-watchdog";
  version = "0.2.0";
  rev = "56bf6aac54deb3863a690f5fc03a2a38e7d9e6ef";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "classic-watchdog";
    rev = version;
    sha256 = "146xp61wd6i1jhifnndsihqx29dwc9snixlniyxp79wcsypg4n5x";
  };

  vendorSha256 = null;

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
