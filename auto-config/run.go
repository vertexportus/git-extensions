package auto_config

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui/list"
	"git_extensions/shared/tui/spinner"
	spinnerBase "github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/spf13/cobra"
)

var argGpgUse bool
var argGpgSign bool
var argName string
var argEmail string
var forceYes bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-auto-config",
	Short: "Helper to auto-config local git repo",
	Long:  `Helper to auto-config local git repo`,
	Args: func(cmd *cobra.Command, args []string) error {
		if argGpgSign && !argGpgUse {
			return fmt.Errorf("sign option requires gpg")
		}
		if !argGpgUse && (argName == "" || argEmail == "") {
			return fmt.Errorf("name and email required if not using GPG source")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		// just set config if not using GPG
		var name string
		var email string
		var gpgKey string
		var sign bool

		if argGpgUse {
			var cancel bool
			cancel, name, email, gpgKey, sign = chooseGpgKey()
			if cancel {
				return
			}
			if argName != "" {
				name = argName
			}
			if argEmail != "" {
				email = argEmail
			}
		} else {
			name = argName
			email = argEmail
		}

		if forceYes || confirm(name, email, sign) {
			err := gitConfig(name, email, gpgKey, sign)
			errors.HandleError(err)
		} else {
			fmt.Println("Aborted")
		}
	},
}

func init() {
	rootCmd.Flags().BoolVarP(
		&argGpgUse,
		"gpg",
		"g",
		false,
		"List GPG keys to use (also selects name and email)")
	rootCmd.Flags().BoolVarP(
		&argGpgSign,
		"sign",
		"s",
		false,
		"Configure git to auto-sign commits")
	rootCmd.Flags().StringVarP(
		&argName,
		"name",
		"n",
		"",
		"Name to configure (overrides name from selected gpg key)")
	rootCmd.Flags().StringVarP(
		&argEmail,
		"email",
		"e",
		"",
		"Email to configure (overrides email from selected gpg key)")
	rootCmd.Flags().BoolVarP(
		&forceYes,
		"yes",
		"y",
		false,
		"Force yes to prompts")
}

var gpgKeys []GpgKey

func Run() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	errors.HandleError(err)
}

var (
	confirmLabelStyle = lipgloss.NewStyle().PaddingLeft(2).Width(24).Foreground(lipgloss.Color("107"))
	confirmValueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("170"))
)

func confirm(name string, email string, sign bool) bool {
	fmt.Println("Configuring git with:")
	fmt.Println(
		lipgloss.JoinHorizontal(0,
			confirmLabelStyle.Render("Name:"),
			confirmValueStyle.Render(name)))
	fmt.Println(
		lipgloss.JoinHorizontal(0,
			confirmLabelStyle.Render("Email:"),
			confirmValueStyle.Render(email)))
	fmt.Println(
		lipgloss.JoinHorizontal(0,
			confirmLabelStyle.Render("Commit auto-signing:"),
			confirmValueStyle.Render(fmt.Sprintf("%t", sign))))
	fmt.Println("")

	input := confirmation.New("Confirm new settings?", confirmation.Yes)
	confirmed, err := input.RunPrompt()
	errors.HandleError(err)
	return confirmed
}

func chooseGpgKey() (bool, string, string, string, bool) {
	// load GPG Keys, showing a spinner animation
	spinner.Show(&spinner.Config{
		Label:   "Loading GPG keys...",
		Spinner: spinnerBase.MiniDot,
		Cmd:     getGpgKeys,
	})

	// finished animation, means gpgKey retrieval finished
	if len(gpgKeys) == 0 {
		fmt.Println("No GPG keys found")
		return true, "", "", "", false
	}

	// render list menu to choose GPG key, quit if none chosen
	gpgKeyListItem, err := list.Choose[GpgKey](
		gpgKeys,
		&list.Config{Title: "Select GPG key", SuppressQuitText: true})
	errors.HandleError(err)
	if gpgKeyListItem == nil {
		return true, "", "", "", false
	}

	gpgKey := gpgKeyListItem.(GpgKey)
	return false, gpgKey.Name, gpgKey.Email, gpgKey.Key, argGpgSign
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
