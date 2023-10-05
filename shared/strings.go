package shared

import "strings"

func TrimExecOutput(output []byte) string {
	return strings.Trim(string(output), "\n")
}
