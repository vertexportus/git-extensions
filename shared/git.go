package shared

import "os/exec"

func GitUpdateConfig(config string, value string) error {
	cmd := exec.Command("git", "config", config, value)
	_, err := cmd.Output()
	return err
}

func GitCurrentBranch() (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return TrimExecOutput(output), nil
}
