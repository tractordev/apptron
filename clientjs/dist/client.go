// dist embeds the final client.js file of clientjs
package dist

import (
	_ "embed"
)

//go:embed client.js
var ClientJS []byte
