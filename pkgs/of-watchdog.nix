{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "of-watchdog";
  version = "0.9.3";
  rev = "4c5f48ffd3f0b566d8ecbf5e359d2320ef10eb60";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "of-watchdog";
    rev = version;
    sha256 = "sha256-t4I+xBg66f9m1c2C9sAgBM6wqM22CZTUsDSLwC/Jerk=";
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
