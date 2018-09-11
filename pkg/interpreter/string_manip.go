package interpret

import (
	"regexp"
	"strings"
)

var spaceStripper = regexp.MustCompile(`\s+`)

func startsWithVowel(s string) bool {
	return strings.HasPrefix(s, "a") ||
		strings.HasPrefix(s, "e") ||
		strings.HasPrefix(s, "i") ||
		strings.HasPrefix(s, "o") ||
		strings.HasPrefix(s, "u")
}

func prepareName(name string) string {
	return defaultStringReplacer.Replace(strings.Replace(strings.ToLower(name), `"`, ``, -1))
}

func clean(s string) string {
	return spaceStripper.ReplaceAllString(s, " ")
}

var defaultStringReplacer = strings.NewReplacer(
	// basic things
	"  ", " ",
	"fmt", "format",
	"sprintf", "sprint f",
	// common initialisms
	"api", "a p i",
	"url", "you are ell",
	"uri", "you are eye",
	// data types
	"[]", "slice of",
	"map[", "map of",
	"bool", "boolean",
	"byte", "bite",
	"complex128", "128-bit complex number",
	"complex64", "64-bit complex number",
	"float32", "64-bit floating point number",
	"float64", "64-bit floating point number",
	"int", " integer integer ",
	"int16", " 16-bit integer ",
	"int32", " 32-bit integer ",
	"int64", " 64-bit integer ",
	"int8", " 8-bit integer ",
	"uint", " unsigned integer ",
	"uint16", " unsigned 16-bit integer ",
	"uint32", " unsigned 32-bit integer ",
	"uint64", " unsigned 64-bit integer ",
	"uint8", " unsigned 8-bit integer ",
	"uintptr", " unsigned integer pointer ",
)
