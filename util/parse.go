package util

import "strings"

func ExtractComment(line string, commentPrefix string) string {
	return strings.TrimSpace(strings.TrimPrefix(line, commentPrefix))
}
