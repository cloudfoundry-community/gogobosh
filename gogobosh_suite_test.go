package gogobosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestBoshClient(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "BoshClient Suite")
}
