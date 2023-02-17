package sqlite

import (
	"encoding/json"

	"github.com/fiatjaf/relayer/storage"
	"github.com/nbd-wtf/go-nostr"
	"golang.org/x/exp/slog"
)

func (b *SqliteBackend) SaveEvent(evt *nostr.Event) error {
	// react to different kinds of events
	if evt.Kind == nostr.KindSetMetadata || evt.Kind == nostr.KindContactList || (10000 <= evt.Kind && evt.Kind < 20000) {
		// delete past events from this user
		b.DB.Exec(`DELETE FROM event WHERE pubkey = $1 AND kind = $2`, evt.PubKey, evt.Kind)
	} else if evt.Kind == nostr.KindRecommendServer {
		// delete past recommend_server events equal to this one
		b.DB.Exec(`DELETE FROM event WHERE pubkey = $1 AND kind = $2 AND content = $3`,
			evt.PubKey, evt.Kind, evt.Content)
	}

	// insert
	tagsj, _ := json.Marshal(evt.Tags)
	res, err := b.DB.Exec(`
        INSERT OR IGNORE INTO event (id, pubkey, created_at, kind, tags, content, sig)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `, evt.ID, evt.PubKey, evt.CreatedAt.Unix(), evt.Kind, tagsj, evt.Content, evt.Sig)
	if err != nil {
		return err
	}

	nr, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if nr == 0 {
		return storage.ErrDupEvent
	}

	return nil
}

func (b *SqliteBackend) BeforeSave(evt *nostr.Event) {
	// do nothing
}

func (b *SqliteBackend) AfterSave(evt *nostr.Event) {
	dbEventsStored.Inc()
	dbEventsTotal.Inc()
	slog.Debug("saved event", "pubkey", evt.PubKey, "kind", evt.Kind)
	// delete all but the 100 most recent ones for each key
	//b.DB.Exec(`DELETE FROM event WHERE pubkey = $1 AND kind = $2 AND created_at < (
	//    SELECT created_at FROM event WHERE pubkey = $1
	//    ORDER BY created_at DESC OFFSET 100 LIMIT 1
	//  )`, evt.PubKey, evt.Kind)
}
