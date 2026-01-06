package locales

import "embed"

//go:embed active.*.toml
var LocalsFS embed.FS
