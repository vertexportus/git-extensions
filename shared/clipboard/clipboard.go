package clipboard

import (
	"git_extensions/shared/errors"
	designclip "golang.design/x/clipboard"
)

func Init() {
	err := designclip.Init()
	errors.HandleError(err)
}

func Write(value string) {
	designclip.Write(designclip.FmtText, []byte(value))
}
