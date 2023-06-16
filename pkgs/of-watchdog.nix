{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "of-watchdog";
  version = "0.9.12";
  rev = "5d4cadcaf595f3d3d213e08cdd38a310c5bb3066";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "of-watchdog";
    rev = version;
    sha256 = "sha256-uXICtmVMT7FnJA0RbCekHoDwtlLn+k94zQhl1bHlHq8=";
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
