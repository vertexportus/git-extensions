package search

import (
	"fmt"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui/list"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"strings"
)

var checkout bool
var pull bool

//var merge bool

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

		if checkout {
			cmd := exec.Command("git", "checkout", branch)
			output, err := cmd.Output()
			errors.HandleError(err)
			fmt.Println(string(output))
			if pull {
				fmt.Println(" ___ Pulling...")
				cmd := exec.Command("git", "pull")
				output, err := cmd.Output()
				errors.HandleError(err)
				fmt.Println(string(output))
			}
		}
		fmt.Println(branch)
	},
}

func init() {
	rootCmd.Flags().BoolVarP(
		&checkout,
		"checkout",
		"c",
		false,
		"checkout to selected branch")
	rootCmd.Flags().BoolVarP(
		&pull,
		"pull",
		"p",
		false,
		"pulls selected branch (either after checkout, or before merge)")
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
	if branch == nil {
		fmt.Println("No branch selected")
		os.Exit(1)
	}
	return branch.(string)
}

func pickFromListMenu(branches []string) any {
	branchListItem, err := list.Choose(
		branches,
		&list.Config{Title: "Select branch", SuppressQuitText: true})
	errors.HandleError(err)

	if branchListItem == "" {
		return ""
	}
	return branchListItem
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
