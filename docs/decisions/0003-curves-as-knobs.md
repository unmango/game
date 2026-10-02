# 0003. Curves are the balancing surface

## Status

Accepted

## Context

Incremental games are built on cost and production functions: the nth upgrade costs `base * growth^n`, production scales polynomially, unlock chances follow a logistic curve.
The framework must let games express these without the framework judging whether the result is balanced.

## Decision

A `Curve` proto message carries a `oneof` over four families:

| Family | Form |
| --- | --- |
| linear | `a + b * n` |
| exponential | `base * growth^n` |
| polynomial | `a * n^k + c` |
| logistic | `max / (1 + e^(-k * (n - mid)))` |

The `curve` package implements `Evaluate(curve, n)`, `Cumulative(curve, from, to)` with closed forms where they exist, `Invert(curve, from, budget)`, and `Advance(curve, level, elapsed)`.
Curve parameters are authored by the game.
The framework offers no opinion about which parameters are balanced.

## Consequences

- A game's balance is a set of curve parameters, which can be tuned without touching framework code.
- New families are additive proto changes.
- Composition of curves (a curve whose parameter is another curve) is out of scope until a game needs it.
