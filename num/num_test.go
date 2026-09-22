package num_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/game/num"
)

var _ = Describe("Number", func() {
	Describe("New", func() {
		It("normalizes a large mantissa", func() {
			n := num.New(1234, 0)
			Expect(n.Mantissa).To(BeNumerically("~", 1.234, 1e-12))
			Expect(n.Exponent).To(Equal(int64(3)))
		})

		It("normalizes a small mantissa", func() {
			n := num.New(0.05, 2)
			Expect(n.Mantissa).To(BeNumerically("~", 5, 1e-12))
			Expect(n.Exponent).To(Equal(int64(0)))
		})

		It("keeps a negative sign", func() {
			n := num.New(-250, 1)
			Expect(n.Mantissa).To(BeNumerically("~", -2.5, 1e-12))
			Expect(n.Exponent).To(Equal(int64(3)))
		})

		It("keeps zero at exponent zero", func() {
			Expect(num.New(0, 99)).To(Equal(num.Zero()))
		})

		It("handles an exact power of ten", func() {
			n := num.New(1000, 0)
			Expect(n.Mantissa).To(BeNumerically("~", 1, 1e-12))
			Expect(n.Exponent).To(Equal(int64(3)))
		})
	})

	Describe("Add", func() {
		It("adds numbers with equal exponents", func() {
			s := num.FromInt(150).Add(num.FromInt(250))
			Expect(s.Float64()).To(BeNumerically("~", 400, 1e-9))
		})

		It("aligns exponents", func() {
			s := num.New(1, 6).Add(num.FromInt(1))
			Expect(s.Float64()).To(BeNumerically("~", 1000001, 1e-3))
		})

		It("drops a term beyond float precision", func() {
			big := num.New(1, 100)
			Expect(big.Add(num.One())).To(Equal(big))
		})

		It("treats zero as identity", func() {
			Expect(num.Zero().Add(num.FromInt(7))).To(Equal(num.FromInt(7)))
			Expect(num.FromInt(7).Add(num.Zero())).To(Equal(num.FromInt(7)))
		})

		It("cancels to zero", func() {
			Expect(num.FromInt(5).Sub(num.FromInt(5)).IsZero()).To(BeTrue())
		})
	})

	Describe("Mul and Div", func() {
		It("multiplies across exponents", func() {
			p := num.New(2, 10).Mul(num.New(3, 20))
			Expect(p.Mantissa).To(BeNumerically("~", 6, 1e-12))
			Expect(p.Exponent).To(Equal(int64(30)))
		})

		It("renormalizes after multiplying", func() {
			p := num.New(5, 0).Mul(num.New(5, 0))
			Expect(p.Mantissa).To(BeNumerically("~", 2.5, 1e-12))
			Expect(p.Exponent).To(Equal(int64(1)))
		})

		It("divides", func() {
			q := num.New(1, 10).Div(num.New(4, 0))
			Expect(q.Mantissa).To(BeNumerically("~", 2.5, 1e-12))
			Expect(q.Exponent).To(Equal(int64(9)))
		})

		It("panics on division by zero", func() {
			Expect(func() { num.One().Div(num.Zero()) }).To(Panic())
		})

		It("exceeds float64 range", func() {
			p := num.New(1, 300).Mul(num.New(1, 300))
			Expect(p.Exponent).To(Equal(int64(600)))
			Expect(math.IsInf(p.Float64(), 1)).To(BeTrue())
		})
	})

	Describe("Pow", func() {
		It("raises to an integer power", func() {
			p := num.FromFloat(1.15).Pow(10)
			Expect(p.Float64()).To(BeNumerically("~", math.Pow(1.15, 10), 1e-9))
		})

		It("raises to a fractional power", func() {
			Expect(num.FromInt(16).Sqrt().Float64()).To(BeNumerically("~", 4, 1e-9))
		})

		It("keeps sign for odd powers of negatives", func() {
			Expect(num.FromInt(-2).Pow(3).Float64()).To(BeNumerically("~", -8, 1e-9))
			Expect(num.FromInt(-2).Pow(2).Float64()).To(BeNumerically("~", 4, 1e-9))
		})

		It("handles zero", func() {
			Expect(num.Zero().Pow(3)).To(Equal(num.Zero()))
			Expect(num.Zero().Pow(0)).To(Equal(num.One()))
		})
	})

	Describe("Log10", func() {
		It("sums the mantissa log and exponent", func() {
			Expect(num.New(1, 45).Log10()).To(BeNumerically("~", 45, 1e-12))
			Expect(num.FromInt(1000).Log10()).To(BeNumerically("~", 3, 1e-12))
		})
	})

	Describe("Cmp", func() {
		DescribeTable("orders values",
			func(a, b num.Number, want int) {
				Expect(a.Cmp(b)).To(Equal(want))
			},
			Entry("equal", num.FromInt(5), num.FromInt(5), 0),
			Entry("by exponent", num.New(1, 5), num.New(9, 4), 1),
			Entry("by mantissa", num.New(2, 5), num.New(3, 5), -1),
			Entry("sign", num.FromInt(-1), num.FromInt(1), -1),
			Entry("negative exponent order", num.New(-1, 5), num.New(-1, 4), -1),
			Entry("zero", num.Zero(), num.FromInt(1), -1),
		)

		It("provides Max and Min", func() {
			a, b := num.FromInt(3), num.FromInt(9)
			Expect(a.Max(b)).To(Equal(b))
			Expect(a.Min(b)).To(Equal(a))
		})
	})

	Describe("formatting", func() {
		DescribeTable("String",
			func(n num.Number, want string) {
				Expect(n.String()).To(Equal(want))
			},
			Entry("zero", num.Zero(), "0"),
			Entry("small", num.FromInt(1500), "1500"),
			Entry("scientific", num.New(1.234, 45), "1.23e45"),
			Entry("negative", num.New(-1.234, 45), "-1.23e45"),
		)

		DescribeTable("Short",
			func(n num.Number, want string) {
				Expect(n.Short()).To(Equal(want))
			},
			Entry("zero", num.Zero(), "0"),
			Entry("units", num.FromInt(5), "5.00"),
			Entry("thousands", num.FromInt(12345), "12.3K"),
			Entry("millions", num.New(1.5, 6), "1.50M"),
			Entry("quadrillions", num.New(1.234, 15), "1.23Qa"),
			Entry("beyond table", num.New(1.234, 45), "1.23e45"),
			Entry("fraction", num.FromFloat(0.5), "0.5"),
		)
	})
})
