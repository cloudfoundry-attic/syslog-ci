package syslog_tests

import (
	"fmt"
	"os"
	"os/exec"

	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var blackbox_config_template = `
---
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

func GenerateBlackBoxConfig(options BlackBoxConfigOptions, dest string) {
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
		Transport: "udp",
		Address:   "127.0.0.1:514",
		SourceDir: ingressDir,
	}

	GenerateBlackBoxConfig(bbOptions, resourceDir)

	cmd := exec.Command(
		"blackbox",
		"-config",
		fmt.Sprintf("%s/blackbox_config.yml", resourceDir),
	)

	bbSession, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	// blackbox needs a little time to start up before being ready to tail files
	time.Sleep(time.Second)

	return bbSession
}
