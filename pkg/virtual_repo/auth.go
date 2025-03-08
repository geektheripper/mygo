package virtual_repo

import (
	"fmt"
	"os"

	"github.com/geektheripper/mygo/cmd/utils"
	"github.com/geektheripper/mygo/pkg/hack_ssh"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
)

func (v *VirtualRepo) EnsureAuth() error {
	if v.auth != nil {
		return nil
	}

	rurl, _ := utils.ParseGitRemoteURL(v.remoteURL)

	if rurl.Protocol == "ssh" {
		if _, ok := os.LookupEnv("SSH_AUTH_SOCK"); ok {
			return nil
		}

		key, err := hack_ssh.GetKeyForHost(rurl.Host, "git")
		if err != nil {
			return err
		}

		auth, err := ssh.NewPublicKeysFromFile("git", key, "")
		if err != nil {
			return err
		}

		v.auth = auth

		return nil
	}

	if rurl.Protocol == "http" {
		fmt.Print("TODO: add http auth")
		return nil
	}

	return nil
}
