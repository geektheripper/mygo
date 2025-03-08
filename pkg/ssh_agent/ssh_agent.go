package ssh_agent

import (
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

type SSHAgent struct {
	sock string
	pid  int
}

var sshSockRegex = regexp.MustCompile(`SSH_AUTH_SOCK=(.*?);`)

func NewSSHAgent() (*SSHAgent, error) {
	a := &SSHAgent{}

	cmd := exec.Command("ssh-agent", "-s")
	cmd.Env = os.Environ()
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	for _, line := range strings.Split(string(output), "\n") {
		if matches := sshSockRegex.FindStringSubmatch(line); len(matches) > 1 {
			a.sock = matches[1]
			break
		}
	}

	a.pid = cmd.Process.Pid

	return a, nil
}

func (a *SSHAgent) SetProcessSock() {
	os.Setenv("SSH_AUTH_SOCK", a.sock)
}

func (a *SSHAgent) Env() []string {
	return []string{
		fmt.Sprintf("SSH_AUTH_SOCK=%s", a.sock),
		fmt.Sprintf("SSH_AGENT_PID=%d", a.pid),
	}
}

func (a *SSHAgent) Command(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	cmd.Env = append(os.Environ(), a.Env()...)
	return cmd
}

func (a *SSHAgent) AddIdentity(path string) error {
	return a.Command("ssh-add", path).Run()
}

func (a *SSHAgent) Kill() error {
	return a.Command("ssh-agent", "-k").Run()
}
