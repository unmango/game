{
  buildGoApplication,
  lib,
  version,
}:
buildGoApplication {
  pname = "game";
  inherit version;

  src = lib.cleanSource ../.;
  modules = ./gomod2nix.toml;

  checkPhase = ''
    runHook preCheck
    go test ./...
    runHook postCheck
  '';
}
