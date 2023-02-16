package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
)

type SqliteBackend struct {
	*sqlx.DB
	DatabaseURL string
	cache       *cache.Cache
}
