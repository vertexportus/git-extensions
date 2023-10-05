package shared

import "golang.design/x/clipboard"

func ClipboardInit() {
	err := clipboard.Init()
	if err != nil {
		panic(err)
	}
}

func ClipboardWrite(value string) {
	clipboard.Write(clipboard.FmtText, []byte(value))
}
