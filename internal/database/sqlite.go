package database

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

// InitDB sets up the database file, configures optimizations, and applies schemas.
func InitDB(dbPath string, schemaPath string) (*sql.DB, error) {
	// 1. Open the physical database connection pool
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open sqlite database: %w", err)
	}

	// Configure connection limits optimal for SQLite embedded runtime single-writer constraints
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)

	// 2. Enforce structural configuration optimizations via PRAGMAs
	// PRAGMA foreign_keys = ON is crucial for Cascading Delete tracking rules to take effect!
	pragmas := []string{
		"PRAGMA foreign_keys = ON;",
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
	}
	for _, pragma := range pragmas {
		if _, err := db.Exec(pragma); err != nil {
			db.Close()
			return nil, fmt.Errorf("failed to execute pragma (%s): %w", pragma, err)
		}
	}

	// 3. Read and execute the master structure definition definitions layout file
	if err := applySchema(db, schemaPath); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to apply database schema: %w", err)
	}

	return db, nil
}

// applySchema parses statements out of the standalone raw sql schema file
func applySchema(db *sql.DB, schemaPath string) error {
	schemaBytes, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("unable to read schema source file: %w", err)
	}

	// Split statements by semicolon delimiter to execute sequentially safely
	queries := strings.Split(string(schemaBytes), ";")
	for _, query := range queries {
		trimmed := strings.TrimSpace(query)
		if trimmed == "" {
			continue
		}
		if _, err := db.Exec(trimmed); err != nil {
			return fmt.Errorf("failed executing schema migration chunk: %w \nQuery: %s", err, trimmed)
		}
	}
	return nil
}
