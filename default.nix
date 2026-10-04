{ lib, buildGoModule }:
buildGoModule {
  pname = "fleur-form";
  version = "0.0.1";

  src = ./.;

  vendorHash = "sha256-J2OWuS0VKkfcwcmK9J/CDBhHyLQ7PlLIiG7KtERhNH4=";

  ldflags = [
    "-s"
    "-w"
  ];

  meta = {
    description = "Fleur router extension";
    homepage = "https://github.com/theobori/fleur-form";
    license = lib.licenses.mit;
  };
}
