{
  buildGoApplication,
  lib,
  go,
  version,
}:
buildGoApplication {
  pname = "game";
  inherit go version;

  src = lib.cleanSource ../.;
  modules = ./gomod2nix.toml;

  checkPhase = ''
    runHook preCheck
    # go test exits 1 when the module has no packages yet.
    if [ -n "$(go list ./... 2>/dev/null)" ]; then
      go test ./...
    fi
    runHook postCheck
  '';
}
