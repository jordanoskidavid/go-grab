package utils

import (
	"net/url"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func NormalizeURL(urlStr string) string {
	u, err := url.Parse(urlStr)
	if err != nil {
		return urlStr
	}
	path := strings.TrimRight(u.Path, "/")
	u.Path = path
	return u.String()
}

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

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

//salted hashing password
/*
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
*/
