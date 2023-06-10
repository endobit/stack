// Package store provides a database connection and queries.
package store

import (
	"context"
	"database/sql"
	_ "embed" // nolint:black-imports

	_ "github.com/mattn/go-sqlite3" // nolint:black-imports

	"github.com/endobit/stack/gen/go/db"
	"github.com/endobit/stack/sql/schema"
)

// Connect returns a database connection and queries.
func Connect() (*db.Queries, error) {
	svc, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	if _, err := svc.ExecContext(context.Background(), schema.SQL); err != nil {
		return nil, err
	}

	return db.New(svc), nil
}
