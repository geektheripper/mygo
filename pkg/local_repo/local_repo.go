package local_repo

import (
	"fmt"
	"path/filepath"

	"gopkg.in/ini.v1"
)

type LocalRepo struct {
	path   string
	config *ini.File
}

func LoadLocalRepo(path string) (*LocalRepo, error) {
	config, err := ini.Load(filepath.Join(path, ".git/config"))
	if err != nil {
		return nil, err
	}

	return &LocalRepo{
		path:   path,
		config: config,
	}, nil
}

func (r *LocalRepo) GetRemoteURL(remoteName string) string {
	return r.config.Section(fmt.Sprintf(`remote "%s"`, remoteName)).Key("url").String()
}
