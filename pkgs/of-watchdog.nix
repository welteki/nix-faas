{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "of-watchdog";
  version = "0.9.6";
  rev = "82eb58722012fe5af0c28dc161a3848e5d3ddd4e";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "of-watchdog";
    rev = version;
    sha256 = "sha256-ssBXoev028k/KfMazfXnb34/rR0SFSHmA3t6iXzCflE=";
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
