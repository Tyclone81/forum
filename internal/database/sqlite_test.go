package database

import (
	"os"
	"testing"
)

func TestInitDB_And_Seeding_Lifecycle(t *testing.T) {
	// Create dynamic runtime workspace references
	tempSchemaFile := "test_schema.sql"
	tempDBFile := "test_forum.db"

	// Mock cleanup hooks
	defer os.Remove(tempSchemaFile)
	defer os.Remove(tempDBFile)

	// Build placeholder minimal valid schema mirroring production file properties
	mockSchemaContent := `
	PRAGMA foreign_keys = ON;
	CREATE TABLE IF NOT EXISTS users (id TEXT PRIMARY KEY, username TEXT UNIQUE);
	CREATE TABLE IF NOT EXISTS categories (id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT UNIQUE);
	`
	if err := os.WriteFile(tempSchemaFile, []byte(mockSchemaContent), 0644); err != nil {
		t.Fatalf("Failed setup preparation stage initializing schema stub assets: %v", err)
	}

	// Execute physical driver validation run loop
	db, err := InitDB(tempDBFile, tempSchemaFile)
	if err != nil {
		t.Fatalf("InitDB returned an unhandled critical execution failure path error: %v", err)
	}
	defer db.Close()

	// Assert storage engine responds correctly to structural verification queries
	var name string
	err = db.QueryRow("SELECT name FROM sqlite_master WHERE type='table' AND name='users';").Scan(&name)
	if err != nil {
		t.Errorf("Validation loop failed, tables failed to compile correctly into the target physical file block: %v", err)
	}

	if name != "users" {
		t.Errorf("Expected table mapping missing target matching structure schema configurations token. Received: %s", name)
	}
}
