package sqlite

import (
	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

func (b *SqliteBackend) Init() error {
	var err error
	b.DB, err = sqlx.Connect("sqlite", b.DatabaseURL)
	if err != nil {
		return err
	}

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
	return err
}
