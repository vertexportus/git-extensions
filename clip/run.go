package clip

import (
	"fmt"
	"git_extensions/shared"
	"os"
)

func Run() {
	shared.ClipboardInit()

	branch, err := shared.GitCurrentBranch()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	shared.ClipboardWrite(branch)
}
