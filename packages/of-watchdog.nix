{ buildGoModule, fetchFromGitHub, ... }:

buildGoModule rec {
  pname = "of-watchdog";
  version = "0.8.4";
  rev = "bbd2e96214264d6b87cc97745ee9f604776dd80f";

  src = fetchFromGitHub {
    owner = "openfaas";
    repo = "of-watchdog";
    rev = version;
    sha256 = "19kg0kf0wf04yapcnbyi58qlxrf1wzlckyxvnnyvpym44zvm7m6d";
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
}
