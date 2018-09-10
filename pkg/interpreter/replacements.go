package interpret

import "strings"

var defaultStringReplacer = strings.NewReplacer(
	"  ", " ",
			"fmt", "format",
			"sprintf", "sprint f",
)
