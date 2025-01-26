package storage

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
)

type ImageMeta struct {
	ID          int64
	ShortID     string
	Filename    string
	ContentType string
	Size        int64
	CreatedAt   time.Time
	ExpiresAt   *time.Time
}

type SQLiteStore struct {
	db *sql.DB
}

func NewSQLiteStore(dbPath string) (*SQLiteStore, error) {
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	store := &SQLiteStore{db: db}
	if err := store.migrate(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *SQLiteStore) migrate() error {
	query := `CREATE TABLE IF NOT EXISTS images (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		short_id TEXT UNIQUE NOT NULL,
		filename TEXT NOT NULL,
		content_type TEXT NOT NULL,
		size INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME
	);
	CREATE INDEX IF NOT EXISTS idx_images_short_id ON images(short_id);`
	_, err := s.db.Exec(query)
	return err
}

func (s *SQLiteStore) Save(meta *ImageMeta) error {
	query := `INSERT INTO images (short_id, filename, content_type, size, expires_at) VALUES (?, ?, ?, ?, ?)`
	_, err := s.db.Exec(query, meta.ShortID, meta.Filename, meta.ContentType, meta.Size, meta.ExpiresAt)
	return err
}

func (s *SQLiteStore) GetByShortID(shortID string) (*ImageMeta, error) {
	meta := &ImageMeta{}
	query := `SELECT id, short_id, filename, content_type, size, created_at, expires_at FROM images WHERE short_id = ?`
	err := s.db.QueryRow(query, shortID).Scan(
		&meta.ID, &meta.ShortID, &meta.Filename, &meta.ContentType,
		&meta.Size, &meta.CreatedAt, &meta.ExpiresAt,
	)
	if err != nil {
		return nil, err
	}
	return meta, nil
}

func (s *SQLiteStore) Close() error {
	return s.db.Close()
}
