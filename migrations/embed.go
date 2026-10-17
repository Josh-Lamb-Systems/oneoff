// Package migrations embeds the SQL schema migrations into the binary so the
// server does not depend on a migrations directory being present on disk.
package migrations

import "embed"

// FS holds the up/down SQL migration files.
//
//go:embed *.sql
var FS embed.FS
