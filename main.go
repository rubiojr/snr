package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rubiojr/snr/sqlite"
	"golang.org/x/exp/slog"

	"github.com/fiatjaf/relayer"
	"github.com/kelseyhightower/envconfig"
	"github.com/nbd-wtf/go-nostr"
)

type Relay struct {
	SqliteDatabase string `envconfig:"SQLITE_DATABASE"`

	storage *sqlite.SqliteBackend
}

func (r *Relay) Name() string {
	return "SNR"
}

func (r *Relay) Storage() relayer.Storage {
	return r.storage
}

func (r *Relay) OnInitialized(*relayer.Server) {}

func (r *Relay) Init() error {
	err := envconfig.Process("", r)
	if err != nil {
		return fmt.Errorf("couldn't process envconfig: %w", err)
	}

	// FIXME: feature flag this
	// every hour, delete all very old events
	//go func() {
	//	db := r.Storage().(*sqlite.SqliteBackend)

	//	for {
	//		time.Sleep(60 * time.Minute)
	//		db.DB.Exec(`DELETE FROM event WHERE created_at < $1`, time.Now().AddDate(0, -3, 0).Unix()) // 3 months
	//	}
	//}()

	return nil
}

// FIXME: feature flag this
func (r *Relay) AcceptEvent(evt *nostr.Event) bool {
	// block events that are too large
	// FIXME: optimize this, perhaps sending and optimization upstream
	// to store the json length before unmarshaling
	jsonb, _ := json.Marshal(evt)
	if len(jsonb) > 10000 {
		return false
	}

	return true
}

func main() {
	var d, v bool
	flag.BoolVar(&d, "debug", false, "Debugging enabled")
	flag.BoolVar(&v, "version", false, "Print version")
	flag.Parse()

	if v {
		fmt.Println(version())
		os.Exit(0)
	}

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		log.Println(http.ListenAndServe(":2112", nil))
	}()

	logLevel := new(slog.LevelVar)
	if d {
		logLevel.Set(slog.LevelDebug)
	}
	h := slog.HandlerOptions{Level: logLevel}.NewTextHandler(os.Stderr)
	slog.SetDefault(slog.New(h))

	r := Relay{}
	if err := envconfig.Process("", &r); err != nil {
		log.Fatalf("failed to read from env: %v", err)
		return
	}

	if r.SqliteDatabase == "" {
		log.Print("Using :memory: SQLite")
		r.SqliteDatabase = ":memory:"
	}
	r.storage = &sqlite.SqliteBackend{DatabaseURL: r.SqliteDatabase}
	if err := relayer.Start(&r); err != nil {
		log.Fatalf("server terminated: %v", err)
	}
}
