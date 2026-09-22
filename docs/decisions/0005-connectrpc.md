# 0005. Services are ConnectRPC over protobuf

## Status

Accepted

## Context

The proto module already exists on the Buf Schema Registry and the first game imports it.
Clients include Go binaries and a browser, and the server should be reachable with `curl` during development.
Plain gRPC needs a proxy for browsers; hand written HTTP drops the proto contract.

## Decision

Services are defined in protobuf and served with ConnectRPC.
`buf.gen.yaml` generates `protoc-gen-go` and `protoc-gen-connect-go` output into `gen/`.
The server binary serves all services on one HTTP listener with h2c enabled.

## Consequences

- One codegen path for every client, including the browser.
- Every endpoint accepts JSON with `Content-Type: application/json`, so smoke tests are one `curl` each.
- gRPC clients still work against the same listener.
