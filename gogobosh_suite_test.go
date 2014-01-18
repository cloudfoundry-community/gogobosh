package gogobosh_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestGoGoBosh(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "GoGoBOSH suite")
}
