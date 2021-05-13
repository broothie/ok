package stringhelp

import "regexp"

var (
	Whitespace    = regexp.MustCompile(`\s+`)
	CommaSplitter = regexp.MustCompile(`\s*,\s*`)
	AllWhitespace = regexp.MustCompile(`^\s*$`).MatchString

	DoubleSlashPrefixMatcher = regexp.MustCompile(`^\s*//`)
	OctothorpePrefixMatcher  = regexp.MustCompile(`^\s*#`)
)

func SplitOnWhitespace(s string) []string {
	return Whitespace.Split(s, -1)
}

func SplitOnCommas(s string) []string {
	return CommaSplitter.Split(s, -1)
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
	if !re.MatchString(s) {
		return nil
	}

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
