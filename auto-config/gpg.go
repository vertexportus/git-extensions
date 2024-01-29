package auto_config

import (
	"fmt"
	strs "git_extensions/shared/strings"
	"log"
	"os/exec"
	"regexp"
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
	// regex to match email and key
	re := regexp.MustCompile(`(?m)(\w+)\s+uid.*]\s([\w\s]+)\s<(.*)>`)
	// get all matches
	matches := re.FindAllStringSubmatch(string(output), -1)
	keys := make([]GpgKey, len(matches))
	// iterate over matches
	for i, match := range matches {
		// get email and key
		key := match[1]
		name := match[2]
		email := match[3]
		// add to map
		keys[i] = GpgKey{Name: name, Email: email, Key: key}
	}
	return keys
}

func getGpgExecPath() (string, error) {
	// get gpg binary
	cmd := exec.Command("git", "config", "gpg.program")
	output, err := cmd.Output()
	if err != nil {
		log.Printf("error reading command: git config gpg.program = %v", err)
		log.Printf("fallbacking to gpg as default path")
		return defaultGpgPath, nil
	}
	cleanOutput := strs.TrimExecOutput(output)
	if cleanOutput == "" {
		return defaultGpgPath, nil
	}
	return cleanOutput, nil
}
