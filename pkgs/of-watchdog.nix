{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "of-watchdog";
  version = "0.9.10";
  rev = "eefeb9dd8c979398a46fc0decc3297591362bfab";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "of-watchdog";
    rev = version;
    sha256 = "sha256-xuA0gPpyLmhEzQHkA+P349nE6elfHZuQAg76kRL9wn0=";
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
    ln -s $out/bin/of-watchdog $out/bin/fwatchdog
  '';
}
