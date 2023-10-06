package clip

import (
	"fmt"
	"git_extensions/shared/clipboard"
	"git_extensions/shared/git"
	"github.com/spf13/cobra"
	"os"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "copies current git branch to clipboard",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		branch, err := git.CurrentBranch()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		clipboard.Write(branch)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
