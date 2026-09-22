# Federation

This document lists the inputs and open questions for player-to-player exchange.
Nothing here is implemented.

## Model

Each server is one player.
A game may register endpoints on the server that accept values from other servers.
The framework's contribution is the exchange rate: a pure function of both servers' public parameters.

Because the rate is pure, an exchange can be computed offline by either party and verified by the other on the next connection.
The ledger records the exchange as an ordinary delta.

## Public parameters

Candidates for the rate function, all readable without authentication:

| Parameter | Source |
| --- | --- |
| age | `created_at` from `IdentityService` |
| level | a game defined path in the ledger, exposed on request |
| rank | a game defined ordering across known peers |
| seed distance | a hash distance between root seeds, for flavor rather than balance |

## Open questions

- Whether identity is the ed25519 public key alone or a key plus a human readable handle.
- Whether a server discovers peers, or games pass peer addresses explicitly.
- Whether an exchange is signed by both parties or only by the sender.
- How a game declares which paths are exchangeable and which are private.
- Whether rate inputs are pinned at offer time or evaluated at settle time.
