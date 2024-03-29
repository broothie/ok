package util

import "strings"

func ExtractCommentIfPresent(line string, commentPrefix string) string {
	if strings.HasPrefix(line, commentPrefix) {
		return strings.TrimSpace(strings.TrimPrefix(line, commentPrefix))
	} else {
		return ""
	}
}
