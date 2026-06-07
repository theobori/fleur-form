{ lib, buildGoModule }:
buildGoModule {
  pname = "fleur-form";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-/Co+seC4PnxQ9LaYU4cg6YSKTi+Z1WQj5REz7JI70XA=";

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
