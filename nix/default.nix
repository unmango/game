{
  buildGoApplication,
  lib,
  ginkgo,
  version,
}:
buildGoApplication {
  pname = "game";
  inherit version;

  src = lib.cleanSource ../.;
  modules = ./gomod2nix.toml;

  nativeCheckInputs = [ ginkgo ];

  checkPhase = ''
    ginkgo run ./...
  '';
}
