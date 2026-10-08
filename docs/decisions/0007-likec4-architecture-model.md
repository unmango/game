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

The CLI is run with `npx` at a pinned version through `make arch`, `make arch-check`, and `make arch-build`, since LikeC4 is not packaged in nixpkgs; the dev shell supplies `nodejs`.
CI validates the model, so a broken reference fails the build.

## Consequences

- The model is plain text and changes alongside the code it describes.
- Rendered output is not committed; `make arch` previews it and `make arch-build` builds a static site into `docs/architecture/dist`.
- Validation needs network access to fetch the CLI, which CI has and the nix build sandbox does not, so it runs outside `nix flake check`.
- Consumers' internals appear here only as far as they touch the framework; each game owns the detail of its own model.
