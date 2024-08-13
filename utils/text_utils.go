package utils

import (
	"strings"
)

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

func FormatTextContent(text string) string {
	words := strings.Fields(text)
	var formattedText strings.Builder
	for i, word := range words {
		formattedText.WriteString(word)
		if (i+1)%10 == 0 {
			formattedText.WriteString(".\n") // End the sentence and start a new line
		} else {
			formattedText.WriteString(" ")
		}
	}
	return strings.TrimSpace(formattedText.String()) // Remove any trailing space
}
