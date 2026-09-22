# 0001. The server is a calculator with an identity

## Status

Accepted

## Context

The framework needs to decide how much state a server holds.
One option is a full game state store that games read and write.
Another is a stateless function library that games embed.
Neither fits: a state store makes the framework know about game content, and a library cannot represent a player to other players.

## Decision

A server holds exactly three persistent facts: a root seed, a creation time, and an ed25519 keypair.
Every numeric answer it gives is a pure function of the root seed and the request.
The public key is the server's identity toward other servers.
The seed and the key are independent, so rotating the key never reseeds the player.

The ledger (ADR 0006) is the one planned addition to persistent state, and it stores deltas rather than game state.

## Consequences

- The Go packages `num`, `seed`, and `curve` have no I/O and can be tested with golden values.
- A game can be rebuilt from the seed and the ledger alone.
- Two servers can compute the same exchange rate without talking to each other.
- Adding any other persistent state requires a new record.
