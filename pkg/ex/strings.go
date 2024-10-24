package ex

import "strings"

func SplitStringOnce(str, sep string) (left string, right string, found bool) {
	sepIndex := strings.Index(str, sep)
	switch sepIndex {
	case -1:
		return str, "", false
	default:
		secondStartIndex := sepIndex + len(sep)
		if secondStartIndex < len(str) {
			return str[0:sepIndex], str[secondStartIndex:], true
		} else {
			return str[0:sepIndex], "", true
		}
	}
}

func TrimLines(lines []string, cutset string) []string {
	return Map(lines, func(line string) string {
		return strings.Trim(line, cutset)
	})
}

func NotEmptyLines(lines []string, trimCutset string) []string {
	return Filtered(lines, func(line string) bool {
		if trimCutset != "" {
			line = strings.Trim(line, trimCutset)
		}
		return line != ""
	})
}
