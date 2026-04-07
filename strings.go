package enumify

import (
	"unicode"
	"unicode/utf8"
)

func LowerFirst(s string) string {
	if s == "" {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}

	return string(unicode.ToLower(r)) + s[size:]
}

func UpperFirst(s string) string {
	if s == "" {
		return ""
	}

	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError {
		return s
	}
	return string(unicode.ToUpper(r)) + s[size:]
}
