package toolhelp

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/broothie/ok/stringhelp"
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
				MatchData: stringhelp.NamedRegexpResult(line, taskMatcher),
				Comment:   stringhelp.Whitespace.ReplaceAllString(strings.TrimSpace(comment.String()), " "),
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
