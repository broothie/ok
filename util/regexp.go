package util

import "regexp"

var paramListCommaSplitter = regexp.MustCompile(`\s*,\s*`)

func SplitCommaParamList(argList string) []string {
	return paramListCommaSplitter.Split(argList, -1)
}

// NamedCaptureGroups pulled from https://stackoverflow.com/questions/20750843/using-named-matches-from-go-regex
func NamedCaptureGroups(re *regexp.Regexp, s string) map[string]string {
	match := re.FindStringSubmatch(s)
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && i < len(match) && name != "" {
			result[name] = match[i]
		}
	}

	return result
}
