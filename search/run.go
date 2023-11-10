package search

import (
	"fmt"
	"git_extensions/shared/cmd"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"git_extensions/shared/tui/list"
	"github.com/erikgeiser/promptkit/confirmation"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

var checkout bool
var pull bool
var merge bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "git-search",
	Short: "Search operations for git",
	Long: `Search operations for git

	It also allows for some shorthand operations on top of found results,
	like merges to current, or checkout and pull`,
	Args: func(cmd *cobra.Command, args []string) error {
		// needs at least the search value
		if err := cobra.MaximumNArgs(1)(cmd, args); err != nil {
			return err
		}
		// cannot use both merge and checkout
		if merge && checkout {
			return fmt.Errorf("cannot use both --merge and --checkout")
		}
		return nil
	},
	Run: func(cobraCmd *cobra.Command, args []string) {
		searchValue := ""
		if len(args) > 0 {
			searchValue = args[0]
		}
		branch := Branch(searchValue, false)

		if checkout {
			checkoutAndPull(branch, pull)
		}
		if merge {
			// get current branch
			currentBranch, err := git.CurrentBranch()
			errors.HandleError(err)
			// checkout to selected branch +pull(if requested)
			checkoutAndPull(branch, pull)
			// checkout back to current branch
			checkoutAndPull(currentBranch, false)

			// confirm merge
			input := confirmation.New(fmt.Sprintf("Merge %s to %s?", branch, currentBranch), confirmation.Yes)
			confirm, err := input.RunPrompt()
			errors.HandleError(err)
			if confirm {
				fmt.Println("merging...")
				cmd.ExecHandleError("git", "merge", branch)
			} else {
				fmt.Println("Aborted")
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
	rootCmd.Flags().BoolVarP(
		&merge,
		"merge",
		"m",
		false,
		"merges selected branch to current branch")
}

func Run() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	err := rootCmd.Execute()
	errors.HandleError(err)
}

func Branch(searchValue string, immediateReturnIfSingle bool) string {
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
	if len(branches) == 1 {
		singleBranch := branches[0]
		if immediateReturnIfSingle {
			return singleBranch
		} else {
			if confirmBranch(singleBranch) {
				return singleBranch
			} else {
				fmt.Println("No branch selected")
				os.Exit(1)
			}
		}
	}
	branch := pickFromListMenu(branches)
	if branch == nil {
		fmt.Println("No branch selected")
		os.Exit(1)
	}
	return branch.(string)
}

func confirmBranch(branch string) bool {
	input := confirmation.New(fmt.Sprintf("Is this the branch you're looking for? '%s'", branch), confirmation.Yes)
	confirm, err := input.RunPrompt()
	errors.HandleError(err)
	return confirm
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

func checkoutAndPull(branch string, pull bool) {
	cmd.ExecHandleError("git", "checkout", branch)
	if pull {
		fmt.Println("pulling...")
		cmd.ExecHandleError("git", "pull")
	}
}
