package util

import (
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"testing"
)

func TestIsEmail(t *testing.T) {
	cases := []struct {
		input    string
		expected bool
	}{
		{"test@test.com", true},
		{"weird-looking+email@domain.com", true},
		{"also_an@email.it", true},
		{"g00d_l00k1nG@3m41L.co.uk", true},
		{"wat", false},
		{"", false},
		{"close@but@not@close@enough", false},
		{"@.", false},
	}

	for _, c := range cases {
		output := IsEmail(c.input)
		equals(t, c.expected, output)
	}
}

func TestGenerateSlug(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"some awesome title", "some-awesome-title"},
		{"aNOTHER aWESOME tITLE", "another-awesome-title"},
		{"oh oh _239", "oh-oh-239"},
		{"new café is opening", "new-cafe-is-opening"},
		{"æ", "ae"},
		{"źŹżŹ", "zzzz"},
		{"Hey ThomasBS😎", "hey-thomasbs"},
	}

	for _, c := range cases {
		output := GenerateSlug(c.input)
		equals(t, c.expected, output)
	}
}

func TestGetMD5Hash(t *testing.T) {
	cases := []struct {
		input, expected string
	}{
		{"123456", "e10adc3949ba59abbe56e057f20f883e"},
		{"hello world", "5eb63bbbe01eeed093cb22bb8f5acdc3"},
		{"👍", "0215ac4dab1ecaf71d83f98af5726984"},
	}

	for _, c := range cases {
		output := GetMD5Hash(c.input)
		equals(t, c.expected, output)
	}
}

// TODO: Move this into its own test package or such for reusability
// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}
