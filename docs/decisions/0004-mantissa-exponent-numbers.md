# 0004. Numbers are mantissa and exponent

## Status

Accepted

## Context

Progression has no upper bound, so the numeric type cannot have one either.
`float64` overflows near `1e308`, which some incrementals reach.
Arbitrary precision types are exact but slow, awkward on the wire, and unavailable to browser clients without a library.

## Decision

The `Number` message is `double mantissa` and `int64 exponent`, normalized so that `1 <= |mantissa| < 10` or the number is zero.
The Go package `num` implements add, sub, mul, div, pow, cmp, log10, conversion from and to `float64`, and formatting in both scientific and short suffix forms.
This is the only numeric type used across the wire for game values.

## Consequences

- Precision is about fifteen significant digits, which is enough to display and compare.
- Adding numbers many orders of magnitude apart drops the smaller one, which is the expected behavior for the genre.
- Clients in any language can implement the type in a few dozen lines.
