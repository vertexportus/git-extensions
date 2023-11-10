package clip

import (
	"git_extensions/search"
	"git_extensions/shared/clipboard"
	shell "git_extensions/shared/cmd"
	"git_extensions/shared/errors"
	"git_extensions/shared/git"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var searchValue string

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "copies current git branch to clipboard",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// get current branch
		var branch string
		if searchValue == "" {
			var err error
			branch, err = git.CurrentBranch()
			errors.HandleError(err)
		} else {
			branch = search.Branch(searchValue, true)
		}
		if runtime.GOOS == "linux" {
			output := shell.ExecHandleError("cat", "/proc/sys/kernel/osrelease")
			if strings.Contains(output, "microsoft") {
				shell.Exec("bash", "-c", "echo "+branch+" | CLIP.EXE")
			}
		} else {
			clipboard.Write(branch)
		}
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
	branchCmd.Flags().StringVarP(
		&searchValue,
		"search",
		"s",
		"",
		"searches for branches containing the search value")
}
