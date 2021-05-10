package toolhelp

import "regexp"

var (
	WhitespaceSplitter = regexp.MustCompile(`\s`)
	AllWhitespace      = regexp.MustCompile(`^\s*$`).MatchString
)

func SplitWhitespace(s string) []string {
	return WhitespaceSplitter.Split(s, -1)
}

func NamedRegexpResults(s string, re *regexp.Regexp) []map[string]string {
	matches := re.FindAllStringSubmatch(s, -1)
	results := make([]map[string]string, len(matches))
	for i, match := range matches {
		results[i] = NamedRegexpResultFromMatches(match, re)
	}

	return results
}

func NamedRegexpResult(s string, re *regexp.Regexp) map[string]string {
	return NamedRegexpResultFromMatches(re.FindStringSubmatch(s), re)
}

func NamedRegexpResultFromMatches(match []string, re *regexp.Regexp) map[string]string {
	result := make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}

	return result
}
