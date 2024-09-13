package utils

import (
	"net/url"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func NormalizeURL(urlStr string) string {
	// Parse the URL string
	u, err := url.Parse(urlStr)
	if err != nil {
		// If there's an error parsing, return the original string
		return urlStr
	}
	// Remove any trailing slash from the path
	path := strings.TrimRight(u.Path, "/")
	u.Path = path
	// Return the modified URL string
	return u.String()
}

func RemoveBlankLines(text string) string {
	var result strings.Builder
	// Split the text into lines
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		// Trim spaces from each line
		trimmed := strings.TrimSpace(line)
		// Only append non-empty lines
		if trimmed != "" {
			// Add a new line if there's already content in the result
			if result.Len() > 0 {
				result.WriteString("\n")
			}
			result.WriteString(trimmed)
		}
	}
	// Return the text with blank lines removed
	return result.String()
}

func RemoveExtraSpaces(text string) string {
	// Split the text into individual words, removing excess spaces
	words := strings.Fields(text)
	// Join the words back into a single string, separated by a single space
	return strings.Join(words, " ")
}

func FormatTextContent(text string) string {
	// Split the text into words
	words := strings.Fields(text)
	var formattedText strings.Builder
	// Loop through each word and format as per the 10-word rule
	for i, word := range words {
		formattedText.WriteString(word)
		// After every 10 words, add a period and a new line
		if (i+1)%10 == 0 {
			formattedText.WriteString(".\n")
		} else {
			formattedText.WriteString(" ")
		}
	}
	// Return the formatted text, trimming any trailing space
	return strings.TrimSpace(formattedText.String()) // Remove any trailing space
}

func HashPassword(password string) (string, error) {
	// Hash the password using bcrypt with the default cost
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		// Return an empty string and the error if hashing fails
		return "", err
	}
	// Return the hashed password as a string
	return string(hashedPassword), nil
}
