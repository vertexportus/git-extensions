package auto_config

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui/list"
	blist "github.com/charmbracelet/bubbles/list"
)

func Run() {
	gpgKeys, err := GetGpgKeys()
	errors.HandleError(err)

	items := make([]blist.Item, len(gpgKeys))
	for i, entry := range gpgKeys {
		items[i] = list.NewListItem(fmt.Sprintf("%s <%s>", entry.Name, entry.Email), entry)
	}

	var gpgKeyListItem list.ItemValue
	gpgKeyListItem, err = list.Choose(&list.Config{Title: "Select GPG key to configure", Items: items})
	errors.HandleError(err)

	err = gitConfig(gpgKeyListItem.(GpgKey))
	errors.HandleError(err)
}

//func runSpinner() {
//	spinnerConfig := spinner.Config{
//		Label:   "Loading GPG keys...",
//		Spinner: bspinner.MiniDot,
//	}
//	spinnerRef = spinner.Show(&spinnerConfig)
//}

func gitConfig(gpgKey GpgKey) error {
	if err := git.UpdateConfig("user.name", gpgKey.Name); err != nil {
		return err
	}
	if err := git.UpdateConfig("user.email", gpgKey.Email); err != nil {
		return err
	}
	if err := git.UpdateConfig("user.signingkey", gpgKey.Key); err != nil {
		return err
	}
	if err := git.UpdateConfig("commit.gpgsign", "true"); err != nil {
		return err
	}

	return nil
}
