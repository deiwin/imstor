package imstor_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestImstor(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Imstor Suite")
}
