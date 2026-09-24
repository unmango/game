package convert_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	gamev1alpha1 "github.com/unmango/game/gen/dev/unmango/game/v1alpha1"

	"github.com/unmango/game/convert"
	"github.com/unmango/game/curve"
	"github.com/unmango/game/num"
)

var _ = Describe("Number", func() {
	It("round trips", func() {
		n := num.New(1234, 5)
		Expect(convert.NumberFromProto(convert.NumberToProto(n))).To(Equal(n))
	})

	It("treats nil as zero", func() {
		Expect(convert.NumberFromProto(nil)).To(Equal(num.Zero()))
	})

	It("normalizes on the way in", func() {
		Expect(convert.NumberFromProto(&gamev1alpha1.Number{Mantissa: 150})).To(Equal(num.New(1.5, 2)))
	})
})

var _ = Describe("Curve", func() {
	DescribeTable("round trips every family",
		func(c curve.Curve) {
			p, err := convert.CurveToProto(c)
			Expect(err).NotTo(HaveOccurred())
			back, err := convert.CurveFromProto(p)
			Expect(err).NotTo(HaveOccurred())
			Expect(back).To(Equal(c))
		},
		Entry("linear", curve.Linear{Intercept: num.FromInt(1), Slope: num.FromInt(2)}),
		Entry("exponential", curve.Exponential{Base: num.FromInt(10), Growth: 1.15}),
		Entry("polynomial", curve.Polynomial{Scale: num.FromInt(2), Degree: 1.5, Offset: num.FromInt(3)}),
		Entry("logistic", curve.Logistic{Max: num.FromInt(100), Steepness: 1, Midpoint: 5}),
	)

	It("rejects an empty family", func() {
		_, err := convert.CurveFromProto(&gamev1alpha1.Curve{})
		Expect(err).To(MatchError(convert.ErrNoFamily))
	})

	It("rejects an unknown curve type", func() {
		_, err := convert.CurveToProto(nil)
		Expect(err).To(MatchError(convert.ErrUnknownCurve))
	})
})
