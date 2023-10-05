package clip

import (
	"fmt"
	"git_extensions/shared"
	"github.com/spf13/cobra"
	"os"
)

// branchCmd represents the branch command
var branchCmd = &cobra.Command{
	Use:   "branch",
	Short: "copies current git branch to clipboard",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		branch, err := shared.GitCurrentBranch()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		shared.ClipboardWrite(branch)
	},
}

func init() {
	rootCmd.AddCommand(branchCmd)
}
