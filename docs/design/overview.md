# Design overview

`unmango/game` is the spine for incremental games that grow without bound.
It answers one question in many forms: given this player's seed and this position in the game, what is the number?

## Responsibilities

The framework does three things.

1. **Identity.** A server is one player's root.
   It holds a root seed, a creation time, and a keypair.
   Nothing else about the player is stored here.
2. **Calculation.** Every price, rate, rarity, and scale factor is a pure function of the seed and the request.
   The server offers a small algebra of curves and derived seeds, and games compose them.
3. **Progression.** An append-only ledger records what the player has done, per path.
   This is the mechanism behind the rule that progress is never lost.
   The ledger is designed here and implemented after the calculator.

The framework does not know what a stat, an upgrade, or a city is.
It does not decide what balanced means.
It exposes the knobs and leaves tuning to the game.

## Concepts

### Root seed

An `int64` chosen at random when the server first starts.
It never changes.
Every other number in the player's world derives from it.

### Paths

Every node in a game is named by a slash separated path such as `character/strength` or `city/restaurant/menu/price`.
The sub-seed for a path is a hash of the root seed and the path (see ADR 0002).
Sub-games are deeper paths, so nesting a game inside a game needs no extra machinery.
This is what makes the tree fractal: any node can be treated as a root by the code below it.

### Curves

A curve maps an integer position `n` to a number.
Families: linear, exponential, polynomial, logistic.
Games author curve parameters; the framework evaluates them.
Three operations cover the incremental genre:

| Operation | Question it answers |
| --- | --- |
| `Evaluate(curve, n)` | what does the nth item cost, or produce |
| `Cumulative(curve, from, to)` | what do items `from` through `to - 1` cost together |
| `Invert(curve, budget)` | how many items can this budget buy starting at `from` |

Parameters may be fixed by the game or derived from a sub-seed, which is how rarity and variance enter without the framework knowing what is rare.

### Numbers

All numbers on the wire are mantissa and exponent (see ADR 0004).
There is no upper bound on progression, so there is no upper bound on the type.

### Time

`Advance(rate, level, elapsed)` integrates a rate curve over a duration.
Live ticking and offline catch-up call the same function with different `elapsed`.
A game that closes for a week reopens by asking what would have happened, and the answer is exact rather than simulated.

### Ledger

The ledger is a list of events, each with a path, a timestamp, and a delta.
State is the fold of the ledger.
Because deltas are only ever added, no action can reduce what the player has earned, which is the framework's one hard rule.
Games decide what a delta means at each path.

### Federation

Two servers can exchange values.
The exchange rate is a pure function of both servers' public parameters (age, level, rank), so it can be computed by either side without a round trip.
See `federation.md` for the inputs and the open questions.

## API surface

`CalculatorService`: `Derive`, `Evaluate`, `Cumulative`, `Invert`, `Advance`.
`IdentityService`: `GetIdentity`.
`LedgerService` (planned): `Append`, `Fold`, `Stream`.
`ExchangeService` (planned): `Rate`, `Offer`, `Settle`.

All services are ConnectRPC over HTTP (see ADR 0005).

## What lives where

Abstractions live here.
Content and rules live in the game.
When a game needs something the framework lacks, the question is whether the addition is math or identity.
If it is, it belongs here.
If it is a rule about what the math means, it belongs in the game.
