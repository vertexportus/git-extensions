package auto_config

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui/list"
	"git_extensions/shared/tui/spinner"
	spinnerBase "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

var gpgKeys []GpgKey

func Run() {
	spinner.Show(&spinner.Config{
		Label:   "Loading GPG keys...",
		Spinner: spinnerBase.MiniDot,
		Cmd:     getGpgKeys,
	})

	if len(gpgKeys) == 0 {
		fmt.Println("No GPG keys found")
		return
	}

	gpgKeyListItem, err := list.Choose[GpgKey](gpgKeys, &list.Config{Title: "Select GPG key"})
	errors.HandleError(err)
	if gpgKeyListItem == nil {
		return
	}

	err = gitConfig(gpgKeyListItem.(GpgKey))
	errors.HandleError(err)
}

func getGpgKeys(channel chan<- tea.Cmd) {
	var err error
	gpgKeys, err = GetGpgKeys()
	errors.HandleError(err)
	close(channel)
}

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
