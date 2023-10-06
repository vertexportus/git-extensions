package git

import (
	"fmt"
	strs "git_extensions/shared/strings"
	"os/exec"
	"strings"
)

func UpdateConfig(config string, value string) error {
	cmd := exec.Command("git", "config", config, value)
	_, err := cmd.Output()
	return err
}

func CurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strs.TrimExecOutput(output), nil
}

func Branches(remote bool, track bool) ([]string, error) {
	if remote && track {
		return nil, fmt.Errorf("cannot specify both remote and track")
	}

	var format string
	if track {
		format = "%(refname:short)|%(upstream:trackshort)"
	} else {
		format = "%(refname:short)"
	}

	var cmd *exec.Cmd
	if remote {
		cmd = exec.Command("git", "branch", "--remote", "--format", format)
	} else {
		cmd = exec.Command("git", "branch", "--format", format)
	}
	output, err := cmd.Output()
	fmt.Sprintln(output)
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strs.TrimExecOutput(output), "\n")

	return branches, nil
}
