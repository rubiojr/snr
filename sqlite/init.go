package sqlite

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	_ "modernc.org/sqlite"
)

const driverOpts = "?_pragma=foreign_keys(1)&_pragma=journal_mode(wal)&cache=shared"

func (b *SqliteBackend) Init() error {
	var err error

	b.cache = cache.New(30*time.Minute, 60*time.Minute)
	b.DB, err = sqlx.Connect("sqlite", b.DatabaseURL+driverOpts)
	if err != nil {
		return err
	}

	b.cache.OnEvicted(func(key string, pq any) {
		prep := pq.(*sql.Stmt)
		prep.Close()
		stmtCacheEvicted.Inc()
	})

	_, err = b.DB.Exec(`
	CREATE TABLE IF NOT EXISTS event(
	  'id' text NOT NULL PRIMARY KEY,
	  'pubkey' text NOT NULL,
	  'created_at' integer NOT NULL,
	  'kind' integer NOT NULL,
	  'tags' text NOT NULL,
	  'content' text NOT NULL,
	  'sig' text NOT NULL
	);
	CREATE UNIQUE INDEX IF NOT EXISTS ididx ON event (id);
	CREATE INDEX IF NOT EXISTS pubkeyprefix ON event (pubkey);
	CREATE INDEX IF NOT EXISTS timeidx ON event (created_at DESC);
	CREATE INDEX IF NOT EXISTS kindidx ON event (kind);
	    `)

	var count uint64
	err = b.DB.QueryRow("SELECT COUNT(*) FROM event").Scan(&count)
	if err != nil {
		return err
	}
	dbEventsTotal.Add(float64(count))

	return nil
}
