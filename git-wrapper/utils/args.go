package utils

import "strings"

func IsFlag(arg string) bool {
	return strings.HasPrefix(arg, "-")
}
