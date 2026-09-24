# 0006. Time is an input and progress is a ledger

## Status

Accepted

## Context

Games run in the foreground, in the background, and not at all.
A player who returns after a week expects the game to have continued.
Simulating the missed ticks is slow and drifts from the live result.
Separately, the framework's one hard rule is that progress is never lost, and it needs a mechanism rather than a convention.

## Decision

`Advance(rate, level, elapsed)` returns the progress accumulated by a rate curve over a duration, using the closed form integral where one exists.
Live ticking calls it with a small `elapsed`; offline catch-up calls it with a large one.
The result is the same function either way.

Progress is stored as an append-only ledger of events, each with a path, a timestamp, and a positive delta.
State is the fold of the ledger.
The ledger is implemented after the calculator, as `LedgerService` with `Append`, `Fold`, and `Stream`.

## Consequences

- Offline progress is exact and instant.
- No operation can reduce an earned value, because the ledger has no negative deltas and no deletes.
- Games that need spending model it as a separate path (spent) rather than a subtraction from earned.
- The ledger's storage format is a later decision; the interface is fixed here.
