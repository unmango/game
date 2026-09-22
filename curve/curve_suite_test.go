package curve_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestCurve(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Curve Suite")
}
