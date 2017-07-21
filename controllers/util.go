package controllers

import (
	"regexp"
)

func IsEmail(email string) bool {
	const email_regex = `^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`
	if m, _ := regexp.MatchString(email_regex, email); !m {
		return false
	}

	return true
}
