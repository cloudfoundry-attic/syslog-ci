package syslog_tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"math/rand"
	"testing"

	"github.com/onsi/ginkgo/config"
)

func TestSyslogTests(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Syslog Tests Suite")
}

var _ = BeforeSuite(func() {
	rand.Seed(config.GinkgoConfig.RandomSeed)
})
