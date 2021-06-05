package util

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"
)

type TaskMatcher func(s string) bool

type RawTask struct {
	Line      int
	MatchData map[string]string
	Comment   string
}

func Scan(r io.Reader, taskMatcher, commentPrefixMatcher *regexp.Regexp) []RawTask {
	var rawTasks []RawTask
	scanner := bufio.NewScanner(r)
	lineCounter := 1
	comment := new(strings.Builder)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if taskMatcher.MatchString(scanner.Text()) {
			rawTasks = append(rawTasks, RawTask{
				Line:      lineCounter,
				MatchData: NamedRegexpResult(line, taskMatcher),
				Comment:   Whitespace.ReplaceAllString(strings.TrimSpace(comment.String()), " "),
			})
		}

		if commentPrefixMatcher.MatchString(line) {
			comment.WriteString(fmt.Sprintf(" %s", commentPrefixMatcher.ReplaceAllString(line, "")))
		} else {
			comment.Reset()
		}

		lineCounter++
	}

	return rawTasks
}
