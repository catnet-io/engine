package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mendsec/catnet-core/pkg/results"
	_ "modernc.org/sqlite"
)

// ScanStore defines the interface for persisting scan data.
type ScanStore interface {
	SaveReport(target string, report *results.ScanReport) (int64, error)
	GetScans() ([]ScanSummary, error)
	GetReport(scanID int64) (*results.ScanReport, error)
	DeleteScan(scanID int64) error
	Close() error
}

type sqliteStore struct {
	db *sql.DB
}

// ScanSummary represents a lightweight view of a historical scan.
type ScanSummary struct {
	ID         int64  `json:"id"`
	StartTime  string `json:"start_time"`
	EndTime    string `json:"end_time"`
	Target     string `json:"target"`
	TotalHosts int    `json:"total_hosts"`
	AliveHosts int    `json:"alive_hosts"`
}

// NewSQLiteStore initializes a SQLite database at the given path.
func NewSQLiteStore(dbPath string) (ScanStore, error) {
	if dbPath != ":memory:" {
		if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
			return nil, fmt.Errorf("failed to create db directory: %w", err)
		}
	}

	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	store := &sqliteStore{db: db}
	if err := store.initSchema(); err != nil {
		db.Close()
		return nil, err
	}

	return store, nil
}

func (s *sqliteStore) initSchema() error {
	schema := `
	CREATE TABLE IF NOT EXISTS scans (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		start_time DATETIME,
		end_time DATETIME,
		target TEXT,
		total_hosts INTEGER,
		alive_hosts INTEGER
	);

	CREATE TABLE IF NOT EXISTS devices (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		scan_id INTEGER,
		ip TEXT,
		hostname TEXT,
		mac TEXT,
		open_ports TEXT, -- JSON array of ints
		is_alive BOOLEAN,
		FOREIGN KEY(scan_id) REFERENCES scans(id) ON DELETE CASCADE
	);
	`
	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to initialize schema: %w", err)
	}
	return nil
}

func (s *sqliteStore) Close() error {
	return s.db.Close()
}
