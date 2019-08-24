package eval_test

import (
	"testing"

	g "github.com/onsi/ginkgo"
	Ω "github.com/onsi/gomega"
)

// Context is used by both eval and ginkgo, so we have to selectively grab what
// we want from ginkgo here
var Describe = g.Describe
var It = g.It
var XIt = g.XIt

func TestEval(t *testing.T) {
	Ω.RegisterFailHandler(g.Fail)
	g.RunSpecs(t, "Eval Suite")
}
