package deployments

import (
	"embed"
)

//go:embed *.sql
var FS embed.FS
