package util

import (
	"regexp"
	"strings"
)

func IsEmail(email string) bool {
	const emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	if m, _ := regexp.MatchString(emailRegex, email); !m {
		return false
	}

	return true
}

func GenerateSlug(title string) string {
	re := regexp.MustCompile("[^a-z0-9]+")
	return strings.Trim(re.ReplaceAllString(strings.ToLower(title), "-"), "-")
}
