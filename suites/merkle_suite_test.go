package suites

import (
	"flag"
	"testing"

	"merkle-challenge/internal/config"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var (
	cfg     *config.Config
	cfgPath string
	errMsg  = "Error %s %s '%v'"
)

func init() { //nolint:gochecknoinits
	// disabling linting: the init is required to declare custom flags before the tests run
	flag.StringVar(&cfgPath, "c", config.DefaultConfigPath, "Config file path")
}

func TestMerkle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Merkle Suite")
}

var _ = BeforeSuite(func() {
	var err error
	cfg, err = config.Get(cfgPath)
	Expect(err).ShouldNot(HaveOccurred())
	Expect(cfg).ShouldNot(BeNil())
})
