package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type Store struct {
	*sql.DB
}

func NewStore(dataSourceName string) (*Store, error) {
	db, err := sql.Open("sqlite3", dataSourceName)
	if err != nil {
		return nil, err
	}

	query := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        username TEXT NOT NULL UNIQUE,
        salt BLOB NOT NULL,
        signing_key BLOB NOT NULL,
        encryption_key BLOB NOT NULL,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	if _, err := db.Exec(query); err != nil {
		return nil, err
	}

	log.Println("Database initialized successfully.")
	return &Store{db}, nil
}

func (s *Store) UserExists(username string) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM users WHERE username = ?)`
	err := s.QueryRow(query, username).Scan(&exists)
	return exists, err
}

func (s *Store) AddUser(username string, salt, signingKey, encryptionKey []byte) error {
	query := `INSERT INTO users (username, salt, signing_key, encryption_key) VALUES (?, ?, ?, ?)`
	_, err := s.Exec(query, username, salt, signingKey, encryptionKey)
	return err
}