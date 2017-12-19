package syslog_tests

import (
	"fmt"
	"os"
	"os/exec"

	"path"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var blackbox_config_template = `---
hostname: local
syslog:
  destination:
    transport: %s
    address: %s
  source_dir: %s
`

type BlackBoxConfigOptions struct {
	Transport string
	Address   string
	SourceDir string
}

func generateBlackBoxConfig(options BlackBoxConfigOptions, dest string) {
	err := os.MkdirAll(dest, 0777)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	file, err := os.Create(dest + "/blackbox_config.yml")
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	defer file.Close()

	n, err := file.Write([]byte(
		fmt.Sprintf(
			blackbox_config_template,
			options.Transport,
			options.Address,
			options.SourceDir,
		)))
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	ExpectWithOffset(1, n).To(BeNumerically(">", 0))
}

func StartBlackbox(ingressDir, resourceDir string) *gexec.Session {
	bbOptions := BlackBoxConfigOptions{
		Transport: "tcp",
		Address:   "127.0.0.1:515",
		SourceDir: path.Dir(ingressDir),
	}

	// TODO - currently we assume that blackbox is installed and on PATH

	generateBlackBoxConfig(bbOptions, resourceDir)

	cmd := exec.Command(
		"blackbox",
		"-config",
		fmt.Sprintf("%s/blackbox_config.yml", resourceDir),
	)

	bbSession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	return bbSession
}
