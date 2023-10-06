package search

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"github.com/spf13/cobra"
	"strings"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-search",
	Short: "Search operations for git",
	Long: `Search operations for git

	It also allows for some shorthand operations on top of found results,
	like merges to current, or checkout and pull`,
	Args: func(cmd *cobra.Command, args []string) error {
		if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
			return err
		}
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		allBranches, err := git.Branches(false, false)
		errors.HandleError(err)
		branches := filterBySearchValue(allBranches, args[0])
		for _, branch := range branches {
			fmt.Println(branch)
		}
	},
}

func Run() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	errors.HandleError(err)
}

func filterBySearchValue(branches []string, searchValue string) []string {
	var filteredBranches []string
	for _, branch := range branches {
		if strings.Contains(branch, searchValue) {
			filteredBranches = append(filteredBranches, branch)
		}
	}
	return filteredBranches
}
