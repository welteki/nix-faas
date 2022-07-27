{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "classic-watchdog";
  version = "0.2.1";
  rev = "cd8dc9f4e98049150d8079a74a18cd5a2e311aeb";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "classic-watchdog";
    rev = version;
    sha256 = "sha256-HMF/1Ky/2WrhogvymqxjbVfsfXH/NZQkRYZPTB6cuD4=";
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
