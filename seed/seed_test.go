package seed_test

import (
	"math"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/game/seed"
)

var _ = Describe("Derive", func() {
	// Golden values computed independently with Python's hashlib.
	// A change here is a breaking change for every game on the framework.
	DescribeTable("matches the pinned derivation",
		func(root int64, path string, want uint64) {
			Expect(seed.Derive(root, path)).To(Equal(want))
		},
		Entry("zero root, empty path", int64(0), "", uint64(12634128529936681850)),
		Entry("zero root", int64(0), "character/strength", uint64(10536520006527416346)),
		Entry("root 42", int64(42), "character/strength", uint64(16750182245374500502)),
		Entry("sibling path", int64(42), "character/agility", uint64(10670125824402861673)),
		Entry("negative root", int64(-1), "city/restaurant", uint64(7227589196759610941)),
		Entry("max root", int64(math.MaxInt64), "a", uint64(2242221872930310847)),
	)

	It("is deterministic", func() {
		Expect(seed.Derive(7, "x/y")).To(Equal(seed.Derive(7, "x/y")))
	})

	It("differs across roots and paths", func() {
		Expect(seed.Derive(1, "x")).NotTo(Equal(seed.Derive(2, "x")))
		Expect(seed.Derive(1, "x")).NotTo(Equal(seed.Derive(1, "y")))
	})
})

var _ = Describe("Rand", func() {
	It("yields the same sequence for the same sub-seed", func() {
		a, b := seed.Rand(99), seed.Rand(99)
		for range 10 {
			Expect(a.Uint64()).To(Equal(b.Uint64()))
		}
	})

	It("yields different sequences for different sub-seeds", func() {
		Expect(seed.Rand(1).Uint64()).NotTo(Equal(seed.Rand(2).Uint64()))
	})
})

var _ = Describe("Join", func() {
	It("joins segments with a slash", func() {
		Expect(seed.Join("character", "strength")).To(Equal("character/strength"))
	})

	It("skips empty segments", func() {
		Expect(seed.Join("", "city", "", "restaurant")).To(Equal("city/restaurant"))
		Expect(seed.Join()).To(Equal(""))
	})
})
