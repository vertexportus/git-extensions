package git

import (
	"fmt"
	"git_extensions/shared/cmd"
	strs "git_extensions/shared/strings"
	"os/exec"
	"strings"
)

func UpdateConfig(config string, value string) error {
	command := exec.Command("git", "config", config, value)
	_, err := command.Output()
	return err
}

func CurrentBranch() (string, error) {
	output, err := cmd.Exec("git", "branch", "--show-current")
	if err != nil {
		return "", err
	}
	return strs.TrimExecOutputStr(output), nil
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

	var command *exec.Cmd
	if remote {
		command = exec.Command("git", "branch", "--remote", "--format", format)
	} else {
		command = exec.Command("git", "branch", "--format", format)
	}
	output, err := command.Output()
	fmt.Sprintln(output)
	if err != nil {
		return nil, err
	}
	branches := strings.Split(strs.TrimExecOutput(output), "\n")

	return branches, nil
}
