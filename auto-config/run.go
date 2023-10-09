package auto_config

import (
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui/list"
)

func Run() {
	gpgKeys, err := GetGpgKeys()
	errors.HandleError(err)

	var gpgKeyListItem any
	gpgKeyListItem, err = list.Choose[GpgKey](gpgKeys, &list.Config{Title: "Select GPG key"})
	errors.HandleError(err)
	if gpgKeyListItem == nil {
		return
	}

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
