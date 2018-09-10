package interpret

import "strings"

var defaultStringReplacer = strings.NewReplacer(
	"  ", " ",
	"fmt", "format",
	"sprintf", "sprint f",
	"api", "a p i",
	"url", "you are ell",
	"uri", "you are eye",
)
