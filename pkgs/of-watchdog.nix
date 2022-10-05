{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "of-watchdog";
  version = "0.9.8";
  rev = "29909ab030e166461ba1ef663a8a293498c10ebb";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "of-watchdog";
    rev = version;
    sha256 = "sha256-HIP54UfwdJ4pO+HwUstWQCe6IpvpP5gGnsYxpRY+YCA=";
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
