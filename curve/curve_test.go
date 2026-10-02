package curve_test

import (
	"math"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/types"

	"github.com/unmango/game/curve"
	"github.com/unmango/game/num"
)

func approx(want float64) types.GomegaMatcher {
	return WithTransform(func(n num.Number) float64 { return n.Float64() },
		BeNumerically("~", want, math.Abs(want)*1e-9+1e-9))
}

var _ = Describe("Linear", func() {
	c := curve.Linear{Intercept: num.FromInt(10), Slope: num.FromInt(3)}

	It("evaluates", func() {
		Expect(c.Evaluate(0)).To(approx(10))
		Expect(c.Evaluate(5)).To(approx(25))
	})

	It("sums in closed form", func() {
		// 10+3*2, 10+3*3, 10+3*4 = 16+19+22
		Expect(curve.Cumulative(c, 2, 5)).To(approx(57))
		Expect(curve.Cumulative(c, 5, 5)).To(approx(0))
	})
})

var _ = Describe("Exponential", func() {
	c := curve.Exponential{Base: num.FromInt(10), Growth: 1.15}

	It("evaluates", func() {
		Expect(c.Evaluate(0)).To(approx(10))
		Expect(c.Evaluate(10)).To(approx(10 * math.Pow(1.15, 10)))
	})

	It("sums in closed form", func() {
		want := 0.0
		for i := 3; i < 8; i++ {
			want += 10 * math.Pow(1.15, float64(i))
		}
		Expect(curve.Cumulative(c, 3, 8)).To(approx(want))
	})

	It("sums a flat growth", func() {
		flat := curve.Exponential{Base: num.FromInt(4), Growth: 1}
		Expect(curve.Cumulative(flat, 0, 5)).To(approx(20))
	})

	It("grows past float64 range", func() {
		Expect(c.Evaluate(10000).Exponent).To(BeNumerically(">", 300))
	})
})

var _ = Describe("Polynomial", func() {
	c := curve.Polynomial{Scale: num.FromInt(2), Degree: 2, Offset: num.FromInt(1)}

	It("evaluates", func() {
		Expect(c.Evaluate(0)).To(approx(1))
		Expect(c.Evaluate(3)).To(approx(19))
	})

	It("sums term by term", func() {
		// 1, 3, 9, 19
		Expect(curve.Cumulative(c, 0, 4)).To(approx(32))
	})
})

var _ = Describe("Logistic", func() {
	c := curve.Logistic{Max: num.FromInt(100), Steepness: 1, Midpoint: 5}

	It("is half of max at the midpoint", func() {
		Expect(c.Evaluate(5)).To(approx(50))
	})

	It("approaches max", func() {
		Expect(c.Evaluate(50).Float64()).To(BeNumerically("~", 100, 1e-6))
	})
})

var _ = Describe("Invert", func() {
	c := curve.Exponential{Base: num.FromInt(10), Growth: 1.15}

	It("returns zero when the first item is unaffordable", func() {
		Expect(curve.Invert(c, 0, num.FromInt(9))).To(Equal(int64(0)))
		Expect(curve.Invert(c, 0, num.Zero())).To(Equal(int64(0)))
	})

	It("returns the count whose cumulative cost fits the budget", func() {
		// 10 + 11.5 + 13.225 = 34.725; adding 15.20875 exceeds 40.
		Expect(curve.Invert(c, 0, num.FromInt(40))).To(Equal(int64(3)))
	})

	It("accepts an exact budget", func() {
		budget := curve.Cumulative(c, 2, 7)
		Expect(curve.Invert(c, 2, budget)).To(Equal(int64(5)))
	})

	It("handles large budgets", func() {
		k := curve.Invert(c, 0, num.New(1, 30))
		Expect(curve.Cumulative(c, 0, k).Cmp(num.New(1, 30))).To(BeNumerically("<=", 0))
		Expect(curve.Cumulative(c, 0, k+1).Cmp(num.New(1, 30))).To(Equal(1))
	})
})

var _ = Describe("Advance", func() {
	rate := curve.Linear{Intercept: num.FromInt(1), Slope: num.FromFloat(0.5)}

	It("multiplies the rate by elapsed seconds", func() {
		Expect(curve.Advance(rate, 4, 10*time.Second)).To(approx(30))
	})

	It("yields nothing for no elapsed time", func() {
		Expect(curve.Advance(rate, 4, 0)).To(Equal(num.Zero()))
		Expect(curve.Advance(rate, 4, -time.Second)).To(Equal(num.Zero()))
	})

	It("is the same for one long interval and many short ones", func() {
		long := curve.Advance(rate, 4, time.Hour)
		short := num.Zero()
		for range 3600 {
			short = short.Add(curve.Advance(rate, 4, time.Second))
		}
		Expect(short.Float64()).To(BeNumerically("~", long.Float64(), 1e-6))
	})
})
