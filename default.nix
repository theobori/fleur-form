{ lib, buildGoModule }:
buildGoModule {
  pname = "fleur-form";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-LsoaUcmYBB0Ktx6cIkhyyMxKFutiiY2hjf0TzMCFg/0=";

  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    description = "";
    homepage = "https://github.com/theobori/fleur-form";
    license = lib.licenses.mit;
  };
}
