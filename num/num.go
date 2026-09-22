// Package num implements a number type with no upper bound, stored as a
// mantissa and a base ten exponent. It is the only numeric type used for game
// values across the framework.
package num

import (
	"math"
	"strconv"
)

// Number is a value stored as mantissa * 10^exponent.
// The zero value is the number zero.
type Number struct {
	Mantissa float64
	Exponent int64
}

// precision is the number of decimal digits a float64 mantissa carries.
// Adding two numbers whose exponents differ by more than this leaves the
// larger unchanged.
const precision = 17

// New returns a normalized Number with the given mantissa and exponent.
func New(mantissa float64, exponent int64) Number {
	return Number{mantissa, exponent}.normalize()
}

// FromFloat converts a float64 to a Number.
func FromFloat(f float64) Number {
	return New(f, 0)
}

// FromInt converts an int64 to a Number.
func FromInt(i int64) Number {
	return New(float64(i), 0)
}

// Zero returns the number zero.
func Zero() Number {
	return Number{}
}

// One returns the number one.
func One() Number {
	return Number{Mantissa: 1}
}

func (n Number) normalize() Number {
	if n.Mantissa == 0 || math.IsNaN(n.Mantissa) || math.IsInf(n.Mantissa, 0) {
		return Number{Mantissa: n.Mantissa}
	}
	shift := int64(math.Floor(math.Log10(math.Abs(n.Mantissa))))
	m := n.Mantissa / math.Pow(10, float64(shift))
	e := n.Exponent + shift
	// Guard against rounding at the boundaries.
	if math.Abs(m) >= 10 {
		m /= 10
		e++
	} else if math.Abs(m) < 1 {
		m *= 10
		e--
	}
	return Number{m, e}
}

// IsZero reports whether n is zero.
func (n Number) IsZero() bool {
	return n.Mantissa == 0
}

// Sign returns -1, 0, or 1 according to the sign of n.
func (n Number) Sign() int {
	switch {
	case n.Mantissa > 0:
		return 1
	case n.Mantissa < 0:
		return -1
	default:
		return 0
	}
}

// Neg returns -n.
func (n Number) Neg() Number {
	return Number{-n.Mantissa, n.Exponent}
}

// Abs returns |n|.
func (n Number) Abs() Number {
	return Number{math.Abs(n.Mantissa), n.Exponent}
}

// Float64 converts n to a float64, saturating to infinity when out of range.
func (n Number) Float64() float64 {
	return n.Mantissa * math.Pow(10, float64(n.Exponent))
}

// Log10 returns the base ten logarithm of |n|.
// It is negative infinity for zero.
func (n Number) Log10() float64 {
	if n.IsZero() {
		return math.Inf(-1)
	}
	return math.Log10(math.Abs(n.Mantissa)) + float64(n.Exponent)
}

// Add returns n + o.
func (n Number) Add(o Number) Number {
	if n.IsZero() {
		return o
	}
	if o.IsZero() {
		return n
	}
	if n.Exponent < o.Exponent {
		n, o = o, n
	}
	diff := n.Exponent - o.Exponent
	if diff > precision {
		return n
	}
	return New(n.Mantissa+o.Mantissa/math.Pow(10, float64(diff)), n.Exponent)
}

// Sub returns n - o.
func (n Number) Sub(o Number) Number {
	return n.Add(o.Neg())
}

// Mul returns n * o.
func (n Number) Mul(o Number) Number {
	return New(n.Mantissa*o.Mantissa, n.Exponent+o.Exponent)
}

// Div returns n / o. Dividing by zero panics, as integer division does.
func (n Number) Div(o Number) Number {
	if o.IsZero() {
		panic("num: division by zero")
	}
	return New(n.Mantissa/o.Mantissa, n.Exponent-o.Exponent)
}

// Pow returns n raised to p.
// A negative n with a non-integer p yields NaN, as math.Pow does.
func (n Number) Pow(p float64) Number {
	if n.IsZero() {
		if p == 0 {
			return One()
		}
		return Zero()
	}
	sign := 1.0
	if n.Sign() < 0 {
		if p != math.Trunc(p) {
			return Number{Mantissa: math.NaN()}
		}
		if int64(p)%2 != 0 {
			sign = -1
		}
	}
	l := n.Log10() * p
	e := math.Floor(l)
	return New(sign*math.Pow(10, l-e), int64(e))
}

// Sqrt returns the square root of n.
func (n Number) Sqrt() Number {
	return n.Pow(0.5)
}

// Cmp compares n and o, returning -1, 0, or 1.
func (n Number) Cmp(o Number) int {
	ns, os := n.Sign(), o.Sign()
	if ns != os {
		return cmp(ns, os)
	}
	if ns == 0 {
		return 0
	}
	if n.Exponent != o.Exponent {
		// For negatives, a larger exponent means a smaller number.
		return ns * cmp(int(n.Exponent-o.Exponent), 0)
	}
	switch {
	case n.Mantissa < o.Mantissa:
		return -1
	case n.Mantissa > o.Mantissa:
		return 1
	default:
		return 0
	}
}

func cmp(a, b int) int {
	switch {
	case a < b:
		return -1
	case a > b:
		return 1
	default:
		return 0
	}
}

// Equal reports whether n and o are the same number.
func (n Number) Equal(o Number) bool {
	return n.Cmp(o) == 0
}

// Less reports whether n < o.
func (n Number) Less(o Number) bool {
	return n.Cmp(o) < 0
}

// Max returns the larger of n and o.
func (n Number) Max(o Number) Number {
	if n.Less(o) {
		return o
	}
	return n
}

// Min returns the smaller of n and o.
func (n Number) Min(o Number) Number {
	if o.Less(n) {
		return o
	}
	return n
}

// String formats n in scientific notation, for example "1.23e45".
// Values below 1e6 in magnitude are formatted as plain decimals.
func (n Number) String() string {
	if n.IsZero() {
		return "0"
	}
	if n.Exponent >= 0 && n.Exponent < 6 {
		return strconv.FormatFloat(n.Float64(), 'f', -1, 64)
	}
	return strconv.FormatFloat(n.Mantissa, 'f', 2, 64) + "e" + strconv.FormatInt(n.Exponent, 10)
}

var suffixes = []string{"", "K", "M", "B", "T", "Qa", "Qi", "Sx", "Sp", "Oc", "No", "Dc"}

// Short formats n with a magnitude suffix and three significant digits, for
// example "12.3K" or "1.23Qa". Beyond the suffix table it falls back to String.
func (n Number) Short() string {
	if n.IsZero() {
		return "0"
	}
	if n.Exponent < 0 {
		return strconv.FormatFloat(n.Float64(), 'g', 3, 64)
	}
	idx := int(n.Exponent / 3)
	if idx >= len(suffixes) {
		return n.String()
	}
	rem := int(n.Exponent % 3)
	v := n.Mantissa * math.Pow(10, float64(rem))
	return strconv.FormatFloat(v, 'f', 2-rem, 64) + suffixes[idx]
}
