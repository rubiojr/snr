package sqlite

import (
	"github.com/jmoiron/sqlx"
	"github.com/patrickmn/go-cache"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	stmtCacheHits = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snr_stmt_cache_hits",
		Help: "Prepared statements cache hits",
	})
	stmtCacheTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snr_stmt_cache_queries",
		Help: "Prepared statements cache queries",
	})
	stmtCachePrepared = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snr_stmt_cache_prepared",
		Help: "Prepared statements cached",
	})
	stmtCacheEvicted = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snr_stmt_cache_evicted",
		Help: "Evicted prepared statements",
	})
	dbEventsStored = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snr_db_events_stored",
		Help: "Database events stored (since last restart)",
	})
	dbEventsTotal = promauto.NewCounter(prometheus.CounterOpts{
		Name: "snr_db_events_total",
		Help: "Total database events stored",
	})
)

type SqliteBackend struct {
	*sqlx.DB
	DatabaseURL string
	cache       *cache.Cache
}
