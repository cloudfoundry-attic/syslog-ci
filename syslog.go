package syslog_tests

import (
	"os/exec"

	. "github.com/onsi/gomega"
)

func RestartSyslog() {
	syslogRestart := exec.Command("sudo", "service", "rsyslog", "restart")
	err := syslogRestart.Run()
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
}
