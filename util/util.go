package util

import (
	"crypto/md5"
	"encoding/hex"
	"regexp"
	"strings"
	"unicode"

	"github.com/rainycape/unidecode"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

func IsEmail(email string) bool {
	const emailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	if m, _ := regexp.MatchString(emailRegex, email); !m {
		return false
	}

	return true
}

func GenerateSlug(title string) string {
	slug := unidecode.Unidecode(title)
	slug = strings.ToLower(slug)
	re := regexp.MustCompile("[^a-z0-9]+")
	slug = re.ReplaceAllString(slug, "-")
	slug = strings.Trim(slug, "-")

	return slug
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func CleanZalgoText(str string) string {
	isMn := func(r rune) bool {
		return unicode.Is(unicode.Mn, r) // Mn: nonspacing marks
	}
	// TODO: Deprecated?
	t := transform.Chain(norm.NFD, transform.RemoveFunc(isMn), norm.NFC)

	result, _, _ := transform.String(t, str)

	return result
}
