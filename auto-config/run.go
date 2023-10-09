package auto_config

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui/list"
	"git_extensions/shared/tui/spinner"
	spinnerBase "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var gpgUse bool
var gpgSign bool
var name string
var email string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-auto-config",
	Short: "Helper to auto-config local git repo",
	Long:  `Helper to auto-config local git repo`,
	Args: func(cmd *cobra.Command, args []string) error {
		if gpgSign && !gpgUse {
			return fmt.Errorf("sign option requires gpg")
		}
		if !gpgUse && (name == "" || email == "") {
			return fmt.Errorf("name and email required if not using GPG source")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		run()
	},
}

func init() {
	rootCmd.Flags().BoolVarP(
		&gpgUse,
		"gpg",
		"g",
		true,
		"(default TRUE) List GPG keys to use (also grabs name and email)")
	rootCmd.Flags().BoolVarP(
		&gpgSign,
		"sign",
		"s",
		false,
		"Configure git to auto-sign commits")
	rootCmd.Flags().StringVarP(
		&name,
		"name",
		"n",
		"",
		"Name to configure (overrides GPG option)")
	rootCmd.Flags().StringVarP(
		&email,
		"email",
		"e",
		"",
		"Email to configure (overrides GPG option)")
}

var gpgKeys []GpgKey

func Run() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	errors.HandleError(err)
}

func run() {
	// just set config if not using GPG
	if !gpgUse {
		err := gitConfig(name, email, "", false)
		errors.HandleError(err)
		return
	}

	// load GPG Keys, showing a spinner anim...
	spinner.Show(&spinner.Config{
		Label:   "Loading GPG keys...",
		Spinner: spinnerBase.MiniDot,
		Cmd:     getGpgKeys,
	})

	// finished animation, means gpgKey retrieval finished
	if len(gpgKeys) == 0 {
		fmt.Println("No GPG keys found")
		return
	}

	// render list menu to choose GPG key, quit if none choosen
	gpgKeyListItem, err := list.Choose[GpgKey](gpgKeys, &list.Config{Title: "Select GPG key"})
	errors.HandleError(err)
	if gpgKeyListItem == nil {
		return
	}

	// configure based on choice
	gpgKey := gpgKeyListItem.(GpgKey)
	err = gitConfig(gpgKey.Name, gpgKey.Email, gpgKey.Key, gpgSign)
	errors.HandleError(err)
}

func getGpgKeys(channel chan<- tea.Cmd) {
	var err error
	gpgKeys, err = GetGpgKeys()
	errors.HandleError(err)
	close(channel)
}

func gitConfig(name string, email string, gpgKey string, sign bool) error {
	if err := gitConfigNameEmail(name, email); err != nil {
		return err
	}
	if err := gitConfigGpg(gpgKey, sign); err != nil {
		return err
	}

	return nil
}

func gitConfigNameEmail(name string, email string) error {
	if err := git.UpdateConfig("user.name", name); err != nil {
		return err
	}
	if err := git.UpdateConfig("user.email", email); err != nil {
		return err
	}

	return nil
}

func gitConfigGpg(gpgKey string, sign bool) error {
	if err := git.UpdateConfig("user.signingkey", gpgKey); err != nil {
		return err
	}

	strSign := "false"
	if sign {
		strSign = "true"
	}
	if err := git.UpdateConfig("commit.gpgsign", strSign); err != nil {
		return err
	}

	return nil
}
