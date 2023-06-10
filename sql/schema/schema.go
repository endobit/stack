// Package schema holds the SQL schema for the stack sqlite database.
package schema

import (
	_ "embed"
)

//go:embed schema.sql

// SQL is the schema for the stack sqlite database.
var SQL string
