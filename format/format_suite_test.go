package format_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestFormat(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Format Suite")
}
