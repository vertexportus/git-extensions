package auto_config

import (
	"fmt"
	strs "git_extensions/shared/strings"
	"log"
	"os/exec"
	"regexp"
	"strings"
)

const defaultGpgPath = "gpg"

type GpgKey struct {
	Name  string
	Email string
	Key   string
}

func (k GpgKey) String() string {
	return fmt.Sprintf("%s <%s>", k.Name, k.Email)
}

func GetGpgKeys() ([]GpgKey, error) {
	// get gpg executable
	gpgExec, err := getGpgExecPath()
	if err != nil {
		return nil, err
	}
	// list keys
	cmd := exec.Command(gpgExec, "--list-secret-keys", "--keyid-format", "LONG")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	return parseGpgKeys(string(output)), nil
}

func parseGpgKeys(output string) []GpgKey {
	lines := strings.Split(output, "\n")
	var keys []GpgKey
	var currentKey string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if match := regexp.MustCompile(`^sec\s+\S+/([A-F0-9]+)`).FindStringSubmatch(line); match != nil {
			currentKey = match[1]
			continue
		}

		if currentKey != "" {
			if match := regexp.MustCompile(`^uid\s+\[\S+\]\s+(.+)\s<(.+)>`).FindStringSubmatch(line); match != nil {
				name := match[1]
				email := match[2]
				keys = append(keys, GpgKey{
					Name:  name,
					Email: email,
					Key:   currentKey,
				})
				currentKey = ""
			}
		}
	}

	return keys
}

func getGpgExecPath() (string, error) {
	// get gpg binary
	cmd := exec.Command("git", "config", "gpg.program")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("gpg.program not configured in git: %v", err)
		log.Printf("Using default GPG path: %s", defaultGpgPath)
		return defaultGpgPath, nil
	}
	cleanOutput := strs.TrimExecOutput(output)
	if cleanOutput == "" {
		return defaultGpgPath, nil
	}
	return cleanOutput, nil
}
