package hack_ssh

import (
	"bufio"
	"fmt"
	"os/exec"
	"regexp"
)

var RegexAttemptKey = []*regexp.Regexp{
	regexp.MustCompile(`Will attempt key: (\S+?)(\s|$)`),
}

var RegexTriedKey = []*regexp.Regexp{
	regexp.MustCompile(`Offering public key: (\S+?)(\s|$)`),
	regexp.MustCompile(`Trying private key: (\S+?)(\s|$)`),
}

var RegexKeyAuthenticated = []*regexp.Regexp{
	regexp.MustCompile(`Authentication succeeded`),
}

type KeyMatch struct {
	AttemptKeys []string
	TriedKeys   []string
	MatchedKey  string
}

func ResolveHost(host string, user ...string) (*KeyMatch, error) {
	if len(user) > 1 {
		return nil, fmt.Errorf("only one user is supported")
	}

	result := &KeyMatch{
		AttemptKeys: make([]string, 0),
		TriedKeys:   make([]string, 0),
		MatchedKey:  "",
	}

	var cmd *exec.Cmd
	if len(user) == 1 {
		cmd = exec.Command("ssh", "-v", user[0]+"@"+host, "exit")
	} else {
		cmd = exec.Command("ssh", "-v", host, "exit")
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		line := scanner.Text()

		for _, regex := range RegexAttemptKey {
			if matches := regex.FindStringSubmatch(line); len(matches) > 1 {
				result.AttemptKeys = append(result.AttemptKeys, matches[1])
			}
		}

		for _, regex := range RegexTriedKey {
			if matches := regex.FindStringSubmatch(line); len(matches) > 1 {
				result.TriedKeys = append(result.TriedKeys, matches[1])
			}
		}

		for _, regex := range RegexKeyAuthenticated {
			if matches := regex.FindStringSubmatch(line); len(matches) > 0 {
				result.MatchedKey = result.TriedKeys[len(result.TriedKeys)-1]
			}
		}
	}

	return result, nil
}

func GetKeyForHost(host string, user ...string) (string, error) {
	match, err := ResolveHost(host, user...)
	if err != nil {
		return "", err
	}
	return match.MatchedKey, nil
}
