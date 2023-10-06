package auto_config

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"github.com/charmbracelet/bubbles/list"
	"os"
)

func Run() {
	gpgKeys, err := GetGpgKeys()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	items := make([]list.Item, len(gpgKeys))
	for i, entry := range gpgKeys {
		items[i] = item{label: fmt.Sprintf("%s <%s>", entry.Name, entry.Email), value: entry}
	}

	var gpgKey GpgKey
	gpgKey, err = ChooseFromList(items, "Select GPG key to configure")
	errors.HandleError(err)

	err = gitConfig(gpgKey)
	errors.HandleError(err)
}

func gitConfig(gpgKey GpgKey) error {
	var err error

	err = git.UpdateConfig("user.name", gpgKey.Name)
	if err != nil {
		return err
	}
	err = git.UpdateConfig("user.email", gpgKey.Email)
	if err != nil {
		return err
	}
	err = git.UpdateConfig("user.signingkey", gpgKey.Key)
	if err != nil {
		return err
	}
	err = git.UpdateConfig("commit.gpgsign", "true")
	if err != nil {
		return err
	}

	return nil
}
