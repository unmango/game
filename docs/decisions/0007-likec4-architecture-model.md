# 0007. Architecture is modeled in LikeC4

## Status

Accepted

## Context

The design documents describe the framework in prose, and the shape of the whole system (the server, its pure packages, the games that consume it, and the planned ledger and federation) is spread across them and the ADRs.
A diagram drawn by hand goes stale the first time a package is added.
The model should live next to the code, be reviewed in pull requests as text, and render views on demand.

## Decision

The architecture is a LikeC4 model in `docs/architecture/`.
`spec.c4` declares element kinds and tags, `model.c4` declares elements and relationships, and `views.c4` declares the views.
Designed but unimplemented parts carry the `#planned` tag and render dashed.
The model covers the framework and its consumers, starting with ouranosis, because the framework's boundary is only meaningful against a game.

LikeC4 is not packaged in nixpkgs, so the flake takes it from `github:unmango/pkgs`, which pins the version and patches it to run from the read-only store.
The dev shell puts `likec4` on the path for `make arch`, `make arch-check`, and `make arch-build`.
A flake check runs `likec4 validate` in the build sandbox, so `nix flake check` and CI fail on a broken reference.

### Amendment: publishing

The model is published as a static site on GitHub Pages.
The flake builds it as `packages.architecture` with `lib.likec4` from `github:UnstoppableMango/a2b`, whose `mangopkgs` input follows `unmango-pkgs`, so the site, the check, and the dev shell share one pinned `likec4`.
The build uses a relative base and hash history, so the same output works at the Pages URL, a custom domain, or opened from disk.
The `pages` workflow builds it with Nix and deploys it on every push to `main`; the `likec4/actions` action is not used because it installs `likec4` from npm outside the pin.
`checks.architecture-site` builds the site on every pull request, so a model that validates but fails to render still fails CI.
Elements link to their source directories and the ADRs that shaped them, so the site is a map into the repository.

## Consequences

- The model is plain text and changes alongside the code it describes.
- Rendered output is not committed; `make arch` previews it, `make arch-build` builds a static site into `docs/architecture/dist`, and `nix build .#architecture` builds the published site.
- Publishing needs Pages enabled with "GitHub Actions" as the source in the repository settings.
- The LikeC4 version moves with the `unmango-pkgs` input; `make update` bumps it with the rest of the flake.
- Consumers' internals appear here only as far as they touch the framework; each game owns the detail of its own model.
