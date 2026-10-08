# 0007. A player may bind an atproto DID to the server key

## Status

Accepted

## Context

ADR 0001 makes the ed25519 public key the server's identity toward other servers.
A raw key is enough for verification but not for people: it has no name, and another server cannot find it from anything a player would say out loud.
`federation.md` leaves two questions open that both come from this gap: whether identity is the key alone or a key plus a human readable handle, and whether servers discover peers or games pass addresses.

The AT Protocol already solves naming and discovery.
An account is a DID, a handle resolves to it and back through DNS or `.well-known`, and the account publishes signed records in a public repository that any server can read.
Players of a game built on this framework are likely to have such an account already.

atproto signing keys are P-256 or secp256k1; ed25519 is not supported.
Making the server key an atproto key would mean changing the key type in ADR 0001 and giving every server its own DID, which a player would never see or name.

## Decision

A player may bind one atproto DID to their server.
Binding is optional: an unbound server keeps every capability it has today, and the ed25519 key stays the identity that signs and verifies framework messages.

The binding is two signatures over the same facts, one from each side.

1. The server signs a binding statement with its ed25519 key: the DID, the server's public endpoint, and `created_at`.
2. The player publishes that statement and its signature as a record in their own repository, under the framework's lexicon `dev.unmango.game.server` with record key `self`.
   The record carries the server's public key as a multibase `did:key` style ed25519 value.
   The repository commit, signed by the DID's key, shows the account agreed to the binding; the ed25519 signature inside shows the server agreed to it.

The server persists the bound DID as a fourth fact next to the seed, creation time, and keypair, and `IdentityService.GetIdentity` reports it.
The handle is not stored, because handles change; it is resolved from the DID when needed.

The server never holds atproto credentials and never writes to a repository.
Signing the statement is pure; publishing it is the game client's job, since the client is where the player signs in.
This keeps the server runnable with no network and keeps atproto out of the pure packages.

Verification by a peer is: resolve the handle to a DID, read `dev.unmango.game.server/self` from the DID's repository, check the ed25519 signature against the key in the record and the DID in the statement, then talk to the endpoint and confirm it presents the same key.

This answers the first two open questions in `federation.md`.
Identity is the ed25519 key, optionally named by a DID and its handle.
Peers are found by resolving a handle; games may still pass addresses directly for unbound servers.

## Alternatives considered

- **Switch the server key to P-256 or secp256k1 and give the server its own DID.** Rejected: it reworks ADR 0001, and a per-server DID has no handle a player would recognize.
- **Declare the server as a service entry in the player's DID document.** Cleaner for discovery, but it needs a PLC operation signed with a rotation key, which most hosted accounts cannot make from an app. It remains open for self-hosted and `did:web` players later.
- **A framework-run handle registry.** Rejected: it adds a central service the framework would have to operate, for a problem atproto already solves.

## Consequences

- `identity.json` gains an optional `did` field; existing files load unchanged.
- `IdentityService` gains the DID in `GetIdentityResponse` and an RPC that returns a signed binding statement for a given DID and endpoint.
- The framework publishes lexicons under the `dev.unmango.game` authority, which requires a `_lexicon` DNS record on the domain that authority names. Lexicon schemas become a second public contract beside the proto module.
- Rotating the server key means a new binding record; per ADR 0001 it never reseeds the player.
- Rebinding to a different DID replaces the record and the stored fact. The old account's record should be deleted by its owner, and peers trust only a record whose statement matches the endpoint's current key.
- Unbinding is a client-side record delete plus clearing the stored DID.
- Later federation records (exchange offers and settlements, ledger checkpoints) can live in the same repository and reference this binding.
