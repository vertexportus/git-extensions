package auto_config

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui"
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
		items[i] = tui.NewListItem(fmt.Sprintf("%s <%s>", entry.Name, entry.Email), entry)
	}

	var gpgKeyListItem tui.ListItemValue
	gpgKeyListItem, err = tui.ChooseFromList(items, "Select GPG key to configure")
	errors.HandleError(err)

	err = gitConfig(gpgKeyListItem.(GpgKey))
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
