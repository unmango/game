package num_test

import (
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestNum(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Num Suite")
}
