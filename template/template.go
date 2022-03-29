package template

import "embed"

//go:embed http-conf https-conf
var TemplateFS embed.FS
