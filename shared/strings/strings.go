package strings

import "strings"

func TrimExecOutput(output []byte) string {
	return TrimExecOutputStr(string(output))
}

func TrimExecOutputStr(output string) string {
	return strings.Trim(output, "\n")
}
