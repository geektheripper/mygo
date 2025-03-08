package utils_test

import (
	"testing"

	"github.com/geektheripper/mygo/cmd/utils"
)

func TestValidateGitRemoteURL(t *testing.T) {
	cases := []struct {
		url      string
		expected bool
	}{
		{"git@github.com:user/repo.git", true},
		{"https://github.com/user/repo.git", true},
		{"https://github.com/user/repo.git", true},
		{"git@github.com:user/repo", false},
	}

	for _, c := range cases {
		if ok := utils.ValidateGitRemoteURL(c.url); ok != c.expected {
			t.Errorf("expected %v for %s", c.expected, c.url)
		}
	}
}
