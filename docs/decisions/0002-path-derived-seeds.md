# 0002. Sub-seeds are derived from the root seed and a path

## Status

Accepted

## Context

Games need a source of variance for every node in their tree: stat rolls, item rarity, cost parameters.
Passing a separate seed per node would make games responsible for seed management and would break the property that everything derives from the root.
Sub-games nested inside games need their own root without any extra registration step.

## Decision

Every node is named by a slash separated path.
The sub-seed for a path is the first eight bytes, big endian, of `SHA-256(root_seed_be64 || path)`.
The Go package `seed` exposes `Derive(root int64, path string) uint64` and `Rand(sub uint64)` which returns a PCG generator from `math/rand/v2` seeded with the sub-seed.

Golden test values pin the derivation.
A change to the hash or the encoding is a breaking change.

## Consequences

- A sub-game rooted at `city/restaurant` treats `Derive(root, "city/restaurant/...")` as its own world; nesting is free.
- Path naming is part of a game's public contract with its players, since renaming a path reseeds it.
- SHA-256 is slower than a non-cryptographic hash, and derivation is rarely on a hot path; games cache sub-seeds if they need to.
