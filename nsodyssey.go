package nsodyssey

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type ProcessNamespaces map[string]string

func Namespaces(pid int) (ProcessNamespaces, error) {
	availableNamespaces, err := listAvailableNamespaces(pid)
	if err != nil {
		return nil, err
	}

	processNamespaces := make(map[string]string)

	for _, namespace := range availableNamespaces {
		inode, err := namespaceInode(strconv.Itoa(pid), namespace)
		if err != nil {
			return nil, err
		}
		processNamespaces[namespace] = inode
	}

	return processNamespaces, nil
}

func (pn ProcessNamespaces) Mnt() string {
	return pn["mnt"]
}

func (pn ProcessNamespaces) Net() string {
	return pn["net"]
}

func (pn ProcessNamespaces) User() string {
	return pn["user"]
}

func (pn ProcessNamespaces) IPC() string {
	return pn["ipc"]
}

func (pn ProcessNamespaces) Pid() string {
	return pn["pid"]
}

func listAvailableNamespaces(pid int) ([]string, error) {
	files, err := ioutil.ReadDir(fmt.Sprintf("/proc/%d/ns", pid))
	if err != nil {
		return nil, err
	}

	availableNamespaces := []string{}
	for _, file := range files {
		availableNamespaces = append(availableNamespaces, file.Name())
	}

	return availableNamespaces, nil
}

func namespaceInode(pid, namespace string) (string, error) {
	inodeLinkContent, err := os.Readlink(fmt.Sprintf("/proc/%s/ns/%s", pid, namespace))
	if err != nil {
		return "", err
	}

	inode, err := Inode(inodeLinkContent)
	if err != nil {
		return "", err
	}

	return inode, nil
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
