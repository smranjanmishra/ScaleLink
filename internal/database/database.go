package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// DB holds the database connection
type DB struct {
	*sql.DB
}

// NewConnection creates a new database connection
func NewConnection(databaseURL string) (*DB, error) {
	db, err := sql.Open("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Initialize tables
	if err := initTables(db); err != nil {
		return nil, fmt.Errorf("failed to initialize tables: %w", err)
	}

	log.Println("✅ Database connected successfully")
	return &DB{db}, nil
}

// initTables creates the necessary tables if they don't exist
func initTables(db *sql.DB) error {
	// Create URLs table
	urlsTable := `
	CREATE TABLE IF NOT EXISTS urls (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		short_code VARCHAR(10) UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		title VARCHAR(255),
		description TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		created_by VARCHAR(100),
		is_active BOOLEAN DEFAULT true,
		expires_at TIMESTAMP,
		INDEX idx_short_code (short_code),
		INDEX idx_created_at (created_at)
	);`

	// Create analytics table
	analyticsTable := `
	CREATE TABLE IF NOT EXISTS analytics (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		url_id UUID NOT NULL,
		short_code VARCHAR(10) NOT NULL,
		ip_address INET,
		user_agent TEXT,
		referer TEXT,
		country VARCHAR(100),
		city VARCHAR(100),
		clicked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		INDEX idx_url_id (url_id),
		INDEX idx_short_code (short_code),
		INDEX idx_clicked_at (clicked_at),
		FOREIGN KEY (url_id) REFERENCES urls(id) ON DELETE CASCADE
	);`

	// Create users table (for future authentication)
	usersTable := `
	CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		password_hash VARCHAR(255) NOT NULL,
		name VARCHAR(255),
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		is_active BOOLEAN DEFAULT true,
		INDEX idx_email (email)
	);`

	// Execute table creation
	tables := []string{urlsTable, analyticsTable, usersTable}
	for _, table := range tables {
		if _, err := db.Exec(table); err != nil {
			return fmt.Errorf("failed to create table: %w", err)
		}
	}

	log.Println("✅ Database tables initialized successfully")
	return nil
}

// Close closes the database connection
func (db *DB) Close() error {
	return db.DB.Close()
} 