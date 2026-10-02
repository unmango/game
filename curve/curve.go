// Package curve evaluates the curve families that price and pace an
// incremental game. See ADR 0003.
package curve

import (
	"math"
	"time"

	"github.com/unmango/game/num"
)

// Curve maps a position to a value.
type Curve interface {
	// Evaluate returns the value at position n.
	Evaluate(n int64) num.Number
}

// Summer is implemented by curves with a closed form for Cumulative.
type Summer interface {
	// Sum returns the total of Evaluate(i) for from <= i < to.
	Sum(from, to int64) num.Number
}

// Linear is intercept + slope * n.
type Linear struct {
	Intercept num.Number
	Slope     num.Number
}

// Evaluate implements Curve.
func (c Linear) Evaluate(n int64) num.Number {
	return c.Intercept.Add(c.Slope.Mul(num.FromInt(n)))
}

// Sum implements Summer.
func (c Linear) Sum(from, to int64) num.Number {
	if to <= from {
		return num.Zero()
	}
	count := num.FromInt(to - from)
	// Sum of i for from <= i < to is count * (from + to - 1) / 2.
	indices := count.Mul(num.FromInt(from + to - 1)).Div(num.FromInt(2))
	return c.Intercept.Mul(count).Add(c.Slope.Mul(indices))
}

// Exponential is base * growth^n.
type Exponential struct {
	Base   num.Number
	Growth float64
}

// Evaluate implements Curve.
func (c Exponential) Evaluate(n int64) num.Number {
	return c.Base.Mul(num.FromFloat(c.Growth).Pow(float64(n)))
}

// Sum implements Summer.
func (c Exponential) Sum(from, to int64) num.Number {
	if to <= from {
		return num.Zero()
	}
	if c.Growth == 1 {
		return c.Base.Mul(num.FromInt(to - from))
	}
	g := num.FromFloat(c.Growth)
	// Geometric series: base * (g^to - g^from) / (g - 1).
	return c.Base.Mul(g.Pow(float64(to)).Sub(g.Pow(float64(from)))).Div(g.Sub(num.One()))
}

// Polynomial is scale * n^degree + offset.
type Polynomial struct {
	Scale  num.Number
	Degree float64
	Offset num.Number
}

// Evaluate implements Curve.
func (c Polynomial) Evaluate(n int64) num.Number {
	return c.Scale.Mul(num.FromInt(n).Pow(c.Degree)).Add(c.Offset)
}

// Logistic is max / (1 + e^(-steepness * (n - midpoint))).
type Logistic struct {
	Max       num.Number
	Steepness float64
	Midpoint  float64
}

// Evaluate implements Curve.
func (c Logistic) Evaluate(n int64) num.Number {
	denom := 1 + math.Exp(-c.Steepness*(float64(n)-c.Midpoint))
	return c.Max.Div(num.FromFloat(denom))
}

// Cumulative returns the total of c at positions from <= i < to.
// Curves with a closed form use it; the rest sum term by term.
func Cumulative(c Curve, from, to int64) num.Number {
	if s, ok := c.(Summer); ok {
		return s.Sum(from, to)
	}
	total := num.Zero()
	for i := from; i < to; i++ {
		total = total.Add(c.Evaluate(i))
	}
	return total
}

// Invert returns the largest count k such that Cumulative(c, from, from+k)
// does not exceed budget. It answers "how many can I afford".
// The curve must be non-negative and non-decreasing from position from.
func Invert(c Curve, from int64, budget num.Number) int64 {
	if budget.Sign() <= 0 || c.Evaluate(from).Cmp(budget) > 0 {
		return 0
	}
	// Gallop to find an upper bound, then binary search.
	lo, hi := int64(1), int64(2)
	for Cumulative(c, from, from+hi).Cmp(budget) <= 0 {
		lo = hi
		if hi > math.MaxInt64/2 {
			return hi
		}
		hi *= 2
	}
	for lo+1 < hi {
		mid := lo + (hi-lo)/2
		if Cumulative(c, from, from+mid).Cmp(budget) <= 0 {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo
}

// Advance returns the progress produced by a rate curve at a fixed level
// over a duration. The rate is per second. Live ticking and offline
// catch-up both call this with different elapsed values.
func Advance(rate Curve, level int64, elapsed time.Duration) num.Number {
	if elapsed <= 0 {
		return num.Zero()
	}
	return rate.Evaluate(level).Mul(num.FromFloat(elapsed.Seconds()))
}
