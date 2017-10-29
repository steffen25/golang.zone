package util

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"bytes"

	"github.com/rainycape/unidecode"
	"golang.org/x/text/runes"
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
	b := make([]byte, len(str))
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	_, _, e := t.Transform(b, []byte(str), true)
	if e != nil {
		panic(e)
	}

	b = bytes.Trim(b, "\x00")

	return string(b)
}

func GetRequestScheme(r *http.Request) string {
	isHTTPS := r.Header.Get("X-Forwarded-Proto") == "https"
	if isHTTPS {
		return "https://"
	}

	return "http://"
}
