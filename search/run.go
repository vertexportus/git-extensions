package search

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui"
	"github.com/charmbracelet/bubbles/list"
	"github.com/spf13/cobra"
	"os"
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
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		searchValue := ""
		if len(args) > 0 {
			searchValue = args[0]
		}
		branch := Branch(searchValue)
		fmt.Println(branch)
	},
}

func Run() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	errors.HandleError(err)
}

func Branch(searchValue string) string {
	// get branches
	allBranches, err := git.Branches(false, false)
	errors.HandleError(err)

	// filter branches by searchValue - if any
	var branches []string
	if searchValue == "" {
		branches = allBranches
	} else {
		branches = filterBySearchValue(allBranches, searchValue)
	}

	// do list menu
	branch := pickFromListMenu(branches)
	if branch == "" {
		fmt.Println("No branch selected")
		os.Exit(1)
	}
	return branch
}

func pickFromListMenu(branches []string) string {
	items := make([]list.Item, len(branches))
	for i, branch := range branches {
		items[i] = tui.NewListSimpleItem(branch)
	}
	branchListItem, err := tui.ChooseFromList(
		&tui.ListConfig{Title: "Select branch", Items: items, SuppressQuitText: true})
	errors.HandleError(err)

	if branchListItem == nil {
		return ""
	}
	return branchListItem.(string)
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
