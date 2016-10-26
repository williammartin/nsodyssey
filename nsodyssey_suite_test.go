package nsodyssey_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestNsodyssey(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Nsodyssey Suite")
}
