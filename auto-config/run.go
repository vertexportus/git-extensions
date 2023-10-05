package auto_config

import (
	"fmt"
	"git_extensions/shared"
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
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = gitConfig(gpgKey)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func gitConfig(gpgKey GpgKey) error {
	var err error

	err = shared.GitUpdateConfig("user.name", gpgKey.Name)
	if err != nil {
		return err
	}
	err = shared.GitUpdateConfig("user.email", gpgKey.Email)
	if err != nil {
		return err
	}
	err = shared.GitUpdateConfig("user.signingkey", gpgKey.Key)
	if err != nil {
		return err
	}
	err = shared.GitUpdateConfig("commit.gpgsign", "true")
	if err != nil {
		return err
	}

	return nil
}
