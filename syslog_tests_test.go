package syslog_tests

import (
	"os"

	"bufio"
	"io/ioutil"

	"fmt"
	"time"

	"net"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	testDir     = "/tmp/syslog-test"
	resourceDir = fmt.Sprintf("%s/resources", testDir)
	ingressDir  = fmt.Sprintf("%s/ingress", testDir)
	egressDir   = fmt.Sprintf("%s/egress", testDir)
	testString  = "hello logging world!"
)

var _ = Describe("syslog", func() {

	// TODO - we can't set up and tear down the directories in an automated way yet:
	// our rsyslog.conf sets ownership of the egress directory to syslog:syslog

	//BeforeEach(func() {
	//	err := os.Mkdir(testDir, 0777)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	err = os.Mkdir(resourceDir, 0777)
	//	Expect(err).NotTo(HaveOccurred())
	//
	//	err = os.Mkdir(ingressDir, 0777)
	//	Expect(err).NotTo(HaveOccurred())
	//})

	//AfterEach(func() {
	//	err := os.RemoveAll(testDir)
	//	Expect(err).ToNot(HaveOccurred())
	//})

	Context("when a log is written directly to rsyslog udp listener", func() {

		BeforeEach(func() {
			// TODO create syslog configuration file

			RestartSyslog()
		})

		It("syslog writes received logs to file", func() {
			conn, err := net.Dial("udp", "127.0.0.1:514")
			Expect(err).ToNot(HaveOccurred())

			n, err := conn.Write([]byte(testString))
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(BeNumerically(">", 0))

			time.Sleep(time.Second) // wait for syslog to process the log entry

			outputBytes := GetOutputBytes()

			Expect(len(outputBytes)).ShouldNot(Equal(0))
			Expect(string(outputBytes)).To(ContainSubstring(testString))
		})
	})

	Context("when a log is written to a file watched by blackbox", func() {

		var (
			bbSession *gexec.Session
		)

		BeforeEach(func() {

			// TODO create syslog configuration file

			RestartSyslog()

			// TODO - currently we assume that blackbox is installed and on PATH

			bbSession = StartBlackbox(ingressDir, resourceDir)
		})

		It("is written to a configured Output Module", func() {

			file, err := os.Create(ingressDir + "/watched_by_blackbox.log")
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()

			n, err := file.Write([]byte(testString + "\n"))
			Expect(err).ToNot(HaveOccurred())
			Expect(n).ToNot(Equal(0))

			time.Sleep(time.Second) // wait for syslog to process the log entry

			// TODO - configure the syslog so that the target depends on test setup

			outputBytes := GetOutputBytes()

			Expect(len(outputBytes)).ShouldNot(Equal(0))
			Expect(string(outputBytes)).To(ContainSubstring(testString))
		})

		AfterEach(func() {
			bbSession.Terminate()
		})
	})
})

func GetOutputBytes() []byte {
	outputFile, err := os.Open(egressDir + "/test.log")
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	outputBytes, err := ioutil.ReadAll(bufio.NewReader(outputFile))
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return outputBytes
}
