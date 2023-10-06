package clip

import (
	"fmt"
	"git_extensions/search"
	"git_extensions/shared/clipboard"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"github.com/spf13/cobra"
	"os"
)

var current bool
var searchValue string

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "copies current git branch to clipboard",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// sanity check
		if current && searchValue != "" {
			fmt.Println("Cannot use both --current and --search")
			os.Exit(1)
		}

		// get current branch
		var branch string
		if current {
			var err error
			branch, err = git.CurrentBranch()
			errors.HandleError(err)
		} else {
			branch = search.Branch(searchValue)
		}
		clipboard.Write(branch)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.Flags().BoolVarP(
		&current,
		"current",
		"c",
		false,
		"copies current git branch to clipboard")
	branchCmd.Flags().StringVarP(
		&searchValue,
		"search",
		"s",
		"",
		"searches for branches containing the search value")
}
