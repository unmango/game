# AGENTS.md

Guidance for AI agents working in this repository.

## Overview

`unmango/game` is a framework for incremental games whose progression can grow without bound.
It is a small server that owns a root seed and answers numeric questions: what an upgrade costs, how fast a rate grows, how much progress elapsed while the player was away.
Games built on it (the first is `github.com/UnstoppableMango/ouranosis`) hold the content and the rules.
This repo holds only the math and the identity behind it.

The design is described in `docs/design/` and the decisions behind it in `docs/decisions/`.
Read `docs/design/overview.md` before changing the API surface.

## Commands

All development tasks go through `make`.
The dev shell (`direnv allow` or `nix develop`) supplies `go`, `buf`, `ginkgo`, and `gomod2nix`.

```sh
make            # nix build plus buf build
make test       # go tool ginkgo run -r
make check      # nix flake check plus buf lint
make fmt        # nix fmt (treefmt: gofmt, nixfmt, actionlint)
make tidy       # go mod tidy and regenerate nix/gomod2nix.toml
make update     # nix flake update
```

Run `make tidy` after any change to `go.mod`; the nix build reads `nix/gomod2nix.toml`, not `go.sum`.

## Architecture

```
proto/dev/unmango/game/v1alpha1/   protobuf messages and services, module buf.build/unmango/game
num/                               mantissa and exponent number type
seed/                              path addressed seed derivation
curve/                             curve evaluation, cumulative cost, inversion, time advance
cmd/game/                          the server binary
docs/design/                       design documents
docs/decisions/                    architecture decision records
nix/                               package derivation and gomod2nix lock
```

The Go packages are pure: no I/O, no clock, no global state.
`cmd/game` is the only place that touches the filesystem or the network.
The proto module is the public contract, and consumers import the generated code from the Buf Schema Registry rather than from this module.

Every numeric result is a deterministic function of the root seed and the request.
Do not add server-side state beyond identity without an ADR; the ledger is the one planned exception and has its own record.

Tests use Ginkgo and Gomega.
Golden values for seed derivation are part of the contract; a change to them is a breaking change for every game built on the framework.
