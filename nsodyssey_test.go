package nsodyssey_test

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"

	"github.com/williammartin/nsodyssey"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Nsodyssey", func() {

	Context("when there is a process running", func() {
		var cmd *exec.Cmd

		BeforeEach(func() {
			cmd = exec.Command("sleep", "60")
			Expect(cmd.Start()).To(Succeed())
		})

		Describe("Namespaces", func() {
			var namespaces nsodyssey.ProcessNamespaces

			BeforeEach(func() {
				var err error
				namespaces, err = nsodyssey.Namespaces(cmd.Process.Pid)
				Expect(err).NotTo(HaveOccurred())
			})

			It("loads the mnt namespace", func() {
				Expect(namespaces.Mnt()).To(Equal(namespaceInode(strconv.Itoa(cmd.Process.Pid), "mnt")))
			})

			It("loads the user namespace", func() {
				Expect(namespaces.User()).To(Equal(namespaceInode(strconv.Itoa(cmd.Process.Pid), "user")))
			})

			It("loads the ipc namespace", func() {
				Expect(namespaces.IPC()).To(Equal(namespaceInode(strconv.Itoa(cmd.Process.Pid), "ipc")))
			})

			It("loads the net namespace", func() {
				Expect(namespaces.Net()).To(Equal(namespaceInode(strconv.Itoa(cmd.Process.Pid), "net")))
			})

			It("loads the pid namespace", func() {
				Expect(namespaces.Pid()).To(Equal(namespaceInode(strconv.Itoa(cmd.Process.Pid), "pid")))
			})

			It("makes namespaces accessible via map syntax", func() {
				Expect(namespaces["mnt"]).To(Equal(namespaceInode(strconv.Itoa(cmd.Process.Pid), "mnt")))
			})
		})
	})
})

func namespaceInode(pid, namespace string) string {
	inodeLinkContent, err := os.Readlink(fmt.Sprintf("/proc/%s/ns/%s", pid, namespace))
	Expect(err).NotTo(HaveOccurred())

	inode, err := Inode(inodeLinkContent)
	Expect(err).NotTo(HaveOccurred())

	return inode
}

func Inode(namespace string) (string, error) {
	requiredNamespaceFormat := regexp.MustCompile(`^\w+:\[\d+\]$`)

	if !requiredNamespaceFormat.MatchString(namespace) {
		return "", fmt.Errorf("namespace string '%s' does not match the required format", namespace)
	}

	namespace = strings.Split(namespace, ":")[1]
	namespace = namespace[1:]
	namespace = namespace[:len(namespace)-1]

	return namespace, nil
}
