package sqlite

import (
	"github.com/jmoiron/sqlx"
)

type SqliteBackend struct {
	*sqlx.DB
	DatabaseURL string
}
