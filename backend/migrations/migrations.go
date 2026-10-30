// Package migrations contains the versioned SQL migration files applied by the
// explicit migration runner. The HTTP server never modifies the schema.
package migrations

import "embed"

// FS exposes the versioned SQL migration files (golang-migrate naming
// convention: <version>_<description>.up.sql) stored in migrations/sql.
//
//go:embed all:sql/*.sql
var FS embed.FS
