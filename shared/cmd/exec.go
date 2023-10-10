package cmd

import (
	"git_extensions/shared/errors"
	"os/exec"
)

func ExecHandleError(args ...string) string {
	output, err := Exec(args...)
	errors.HandleError(err)
	return output
}

func Exec(args ...string) (string, error) {
	command := exec.Command(args[0], args[1:]...)
	output, err := command.Output()
	if err != nil {
		return "", err
	}
	strOutput := string(output)
	return strOutput, err
}
