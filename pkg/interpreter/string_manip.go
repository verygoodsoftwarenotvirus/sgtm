package interpret

import(
	"strings"
)

func startsWithVowel(s string) bool {
	return strings.HasPrefix(s, "a") ||
		strings.HasPrefix(s, "e")||
		strings.HasPrefix(s, "i")||
		strings.HasPrefix(s, "o")||
		strings.HasPrefix(s, "u")
}
