package static

import "embed"

//go:embed *.html *.css *.js
var StaticFolder embed.FS
