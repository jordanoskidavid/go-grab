package utils

import "strings"

// Filter functions for removing the blank and extra spaces, all to be tidy and clean
func RemoveBlankLines(text string) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			if result.Len() > 0 {
				result.WriteString("\n")
			}
			result.WriteString(trimmed)
		}
	}
	return result.String()
}

func RemoveExtraSpaces(text string) string {
	words := strings.Fields(text)
	return strings.Join(words, " ")
}
