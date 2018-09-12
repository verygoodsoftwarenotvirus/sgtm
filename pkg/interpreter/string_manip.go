package interpret

import (
	"regexp"
	"strings"
	"unicode"
)

var spaceStripper = regexp.MustCompile(`\s+`)

func startsWithVowel(s string) bool {
	return strings.HasPrefix(s, "a") ||
		strings.HasPrefix(s, "e") ||
		strings.HasPrefix(s, "i") ||
		strings.HasPrefix(s, "o") ||
		strings.HasPrefix(s, "u")
}

func splitByCap(input string) []string {
	var words []string
	l := 0
	for s := input; s != ""; s = s[l:] {
		l = strings.IndexFunc(s[1:], unicode.IsUpper) + 1
		if l <= 0 {
			l = len(s)
		}
		words = append(words, s[:l])
	}
	return words

}

func prepareName(name string) string {
	first := strings.Replace(name, `"`, ``, -1)
	second := splitByCap(first)
	third := strings.Join(second, " ")
	fourth := replace(strings.ToLower(third))

	return fourth
}

func clean(s string) string {
	return spaceStripper.ReplaceAllString(s, " ")
}

func replace(s string) string {
	dataTypeReplacements := map[string]string{
		// data types
		"^bool$":       " boolean ",
		"^byte$":       " bite ",
		"^complex128$": " 128-bit complex number ",
		"^complex64$":  " 64-bit complex number ",
		"^float32$":    " 64-bit floating point number ",
		"^float64$":    " 64-bit floating point number ",
		"^int$":        " integer ",
		"^int16$":      " 16-bit integer ",
		"^int32$":      " 32-bit integer ",
		"^int64$":      " 64-bit integer ",
		"^int8$":       " 8-bit integer ",
		"^uint$":       " unsigned integer ",
		"^uint16$":     " unsigned 16-bit integer ",
		"^uint32$":     " unsigned 32-bit integer ",
		"^uint64$":     " unsigned 64-bit integer ",
		"^uint8$":      " unsigned 8-bit integer ",
		"^uintptr$":    " unsigned integer pointer ",
	}

	for o, n := range dataTypeReplacements {
		s = regexp.MustCompile(o).ReplaceAllString(s, n)
	}

	return defaultStringReplacer.Replace(s)
}

var defaultStringReplacer = strings.NewReplacer(
	// basic things
	"fmt", " format",
	"sprintf", " sprint f ",
	"println", " print line ",
	"printf", " print f ",
	// common initialisms
	"api", " a p i ",
	"url", " you are ell ",
	"uri", " you are eye ",
)
