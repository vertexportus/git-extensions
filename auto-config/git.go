package auto_config

import "os/exec"

func UpdateGitConfig(config string, value string) error {
	cmd := exec.Command("git", "config", config, value)
	_, err := cmd.Output()
	return err
}
