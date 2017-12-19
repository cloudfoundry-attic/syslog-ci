package syslog_tests

import (
	"os"

	"bufio"
	"io/ioutil"

	"fmt"
	"time"

	"net"

	"strconv"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
)

var (
	testDir        = "/tmp/syslog-test"
	resourceDir    = fmt.Sprintf("%s/resources", testDir)
	ingressDir     = fmt.Sprintf("%s/ingress/test", testDir)
	egressFilename = fmt.Sprintf("%s/egress/output_from_syslog.log", testDir)
	basicTestLine  = counterString(100, "*")
	longTestLine   = counterString(1025, "*")
)

var _ = Describe("syslog", func() {

	// TODO - we can't set up and tear down the directories in an automated way yet:
	// our rsyslog.conf sets ownership of the egress directory to syslog:syslog

	Context("when a short log is written directly to rsyslog udp listener", func() {

		BeforeEach(func() {
			CreateFolders()

			GenerateTestOutputConfig(egressFilename)
			RestartSyslog()
		})

		It("syslog writes received logs to file", func() {
			conn, err := net.Dial("udp", "127.0.0.1:514")
			Expect(err).ToNot(HaveOccurred())

			n, err := conn.Write([]byte(basicTestLine + "\n"))
			Expect(err).ToNot(HaveOccurred())
			Expect(n).To(BeNumerically(">", 0))

			time.Sleep(time.Second) // wait for syslog to process the log entry

			outputBytes := GetOutputBytes()
			Expect(len(outputBytes)).ShouldNot(Equal(0))
			Expect(string(outputBytes)).To(ContainSubstring(basicTestLine))
		})

		AfterEach(func() {
			CleanupTestDir()
		})

	})

	Context("when a short log is written to a file watched by blackbox", func() {

		var (
			bbSession *gexec.Session
		)

		BeforeEach(func() {
			CreateFolders()

			GenerateTestOutputConfig(egressFilename)
			RestartSyslog()

			bbSession = StartBlackbox(ingressDir, resourceDir)
		})
		//
		It("is written to a configured Output Module", func() {

			file, err := os.Create(ingressDir + "/watched_by_blackbox.log")
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()

			time.Sleep(5 * time.Second) // wait for syslog to process the log entry

			n, err := file.Write([]byte(basicTestLine + "\n"))
			Expect(err).ToNot(HaveOccurred())
			Expect(n).ToNot(Equal(0))

			time.Sleep(time.Second) // wait for syslog to process the log entry
			// TODO - configure the syslog so that the target depends on test setup

			outputBytes := GetOutputBytes()
			Expect(len(outputBytes)).ShouldNot(Equal(0))
			Expect(string(outputBytes)).To(ContainSubstring(basicTestLine))

		})

		AfterEach(func() {
			bbSession.Terminate()
			CleanupTestDir()
		})
	})

	Context("when a long log is written to a file watched by blackbox", func() {

		var (
			bbSession *gexec.Session
		)

		BeforeEach(func() {
			CreateFolders()

			GenerateTestOutputConfig(egressFilename)
			RestartSyslog()

			bbSession = StartBlackbox(ingressDir, resourceDir)
		})

		It("is written to a configured Output Module", func() {

			file, err := os.Create(ingressDir + "/watched_by_blackbox.log")
			Expect(err).ToNot(HaveOccurred())
			defer file.Close()

			time.Sleep(5 * time.Second)

			n, err := file.Write([]byte(longTestLine + "\n"))
			Expect(err).ToNot(HaveOccurred())
			Expect(n).ToNot(Equal(0))

			time.Sleep(time.Second)

			outputBytes := GetOutputBytes()
			Expect(len(outputBytes)).ShouldNot(Equal(0))
			Expect(string(outputBytes)).To(ContainSubstring(longTestLine))
		})

		AfterEach(func() {
			bbSession.Terminate()
			CleanupTestDir()
		})
	})
})

func CreateFolders() {
	err := os.Mkdir(testDir, 0777)
	Expect(err).NotTo(HaveOccurred())
	err = os.MkdirAll(ingressDir, 0777)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
}

func CleanupTestDir() {
	err := os.RemoveAll(testDir)
	Expect(err).ToNot(HaveOccurred())
}

func GetOutputBytes() []byte {
	outputFile, err := os.Open(egressFilename)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())
	outputBytes, err := ioutil.ReadAll(bufio.NewReader(outputFile))
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	return outputBytes
}

func counterString(l int, s string) string {
	counterstring := ""
	for len(counterstring) < l {
		counterstring += s
		counterstring += strconv.Itoa(len(counterstring))
	}

	return counterstring[:l]
}
