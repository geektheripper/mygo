package utils

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var GitRemoteRegexSSH = regexp.MustCompile(`^git@([a-zA-Z0-9-_\.]+)\:(.*).git$`)

type GitRemote struct {
	URL      string
	Protocol string
	Host     string
	Port     int
}

func ParseGitRemoteURL(remoteURL string) (*GitRemote, error) {
	if !strings.HasSuffix(remoteURL, ".git") {
		return nil, fmt.Errorf("invalid remote URL: %s", remoteURL)
	}

	if gurl, err := url.Parse(remoteURL); err == nil {
		if gurl.Scheme != "http" && gurl.Scheme != "https" {
			return nil, fmt.Errorf("invalid remote URL: %s", remoteURL)
		}

		return &GitRemote{
			URL:      remoteURL,
			Protocol: gurl.Scheme,
			Host:     gurl.Host,
		}, nil
	}

	if matches := GitRemoteRegexSSH.FindStringSubmatch(remoteURL); len(matches) == 3 {
		return &GitRemote{
			URL:      remoteURL,
			Protocol: "ssh",
			Host:     matches[1],
		}, nil
	}

	return nil, fmt.Errorf("invalid remote URL: %s", remoteURL)
}

func ValidateGitRemoteURL(remoteURL string) bool {
	_, err := ParseGitRemoteURL(remoteURL)
	return err == nil
}
