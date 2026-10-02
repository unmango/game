package identity_test

import (
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/unmango/game/identity"
)

var _ = Describe("Identity", func() {
	now := time.Date(2026, 9, 21, 12, 0, 0, 0, time.UTC)

	It("generates distinct identities", func() {
		a, err := identity.Generate(now)
		Expect(err).NotTo(HaveOccurred())
		b, err := identity.Generate(now)
		Expect(err).NotTo(HaveOccurred())
		Expect(a.Seed).NotTo(Equal(b.Seed))
		Expect(a.PublicKey()).NotTo(Equal(b.PublicKey()))
		Expect(a.CreatedAt).To(Equal(now))
	})

	It("round trips through a data directory", func() {
		dir := filepath.Join(GinkgoT().TempDir(), "data")
		first, err := identity.LoadOrGenerate(dir, now)
		Expect(err).NotTo(HaveOccurred())
		second, err := identity.LoadOrGenerate(dir, now.Add(time.Hour))
		Expect(err).NotTo(HaveOccurred())
		Expect(second).To(Equal(first))
	})

	It("reports a missing identity", func() {
		_, err := identity.Load(GinkgoT().TempDir())
		Expect(err).To(MatchError(ContainSubstring("no such file")))
	})
})
