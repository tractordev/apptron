// chrome embeds the builtin chrome options
package chrome

import (
	"embed"
)

//go:embed *.html
var Dir embed.FS
