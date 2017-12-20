package syslog_tests

import (
	"os/exec"

	"fmt"
	"os"

	. "github.com/onsi/gomega"
)

var syslog_config_template = `
$umask 0000
$ModLoad imuxsock
$ModLoad imudp
$ModLoad imtcp
$UDPServerAddress 127.0.0.1
$UDPServerRun 514
$InputTCPServerRun 515

action(
  type="omfile"
  dirCreateMode="0777"
  fileCreateMode="0777"
  dirOwner="syslog"
  dirGroup="syslog"
  fileOwner="syslog"
  fileGroup="syslog"
  file="%s"
)
`

type SyslogConfigOptions struct {
	EgressFile string
	EgressDir  string
}

func GenerateTestOutputConfig(outputFilepath string) {
	file, err := os.Create("/etc/rsyslog.d/00-test.conf")
	// Important Note: this requires the user running the tests
	// to have write access to rsyslog.d.
	ExpectWithOffset(1, err).ToNot(HaveOccurred())

	defer file.Close()

	bytesWritten, err := file.Write([]byte(
		fmt.Sprintf(
			syslog_config_template,
			outputFilepath,
		)))
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	ExpectWithOffset(1, bytesWritten).To(BeNumerically(">", 0))
}

func RestartSyslog() {
	syslogRestart := exec.Command("sudo", "service", "rsyslog", "restart")
	err := syslogRestart.Run()
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
}
