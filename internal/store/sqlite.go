package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

type Store struct {
	db *sql.DB
}

type Post struct {
	ID          int
	Slug        string
	Title       string
	Summary     string
	HTML        string
	RawMD       string
	PublishedAt string
	UpdatedAt   string
	Draft       bool
}

type Tag struct {
	ID   int
	Name string
}

// MustOpen opens a database (SQLite or Turso) based on environment variables
func MustOpen(dbPath string) *Store {
	var db *sql.DB
	var err error

	// Check if Turso credentials are available
	dbURL := os.Getenv("DB_URL")
	dbToken := os.Getenv("DB_AUTH_TOKEN")
	
	if dbURL != "" && dbToken != "" {
		// Use Turso/libSQL  
		dsn := fmt.Sprintf("%s?authToken=%s", dbURL, dbToken)
		db, err = sql.Open("libsql", dsn)
		if err != nil {
			panic(fmt.Sprintf("failed to open Turso database: %v", err))
		}
	} else {
		// Fall back to local SQLite
		db, err = sql.Open("sqlite3", dbPath)
		if err != nil {
			panic(fmt.Sprintf("failed to open SQLite database: %v", err))
		}
	}

	// Test connection
	if err := db.Ping(); err != nil {
		panic(fmt.Sprintf("failed to connect to database: %v", err))
	}

	store := &Store{db: db}
	if err := store.migrate(); err != nil {
		panic(fmt.Sprintf("failed to migrate database: %v", err))
	}

	return store
}

// migrate applies SQL migrations from migrations/ directory
func (s *Store) migrate() error {
	// Create migrations table if it doesn't exist
	_, err := s.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version TEXT PRIMARY KEY,
			applied_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Read migrations from filesystem
	migrations, err := s.loadMigrations()
	if err != nil {
		return fmt.Errorf("failed to load migrations: %w", err)
	}

	// Apply each migration
	for _, migration := range migrations {
		if applied, err := s.isMigrationApplied(migration.Version); err != nil {
			return fmt.Errorf("failed to check migration %s: %w", migration.Version, err)
		} else if applied {
			continue
		}

		// Execute migration
		if _, err := s.db.Exec(migration.SQL); err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", migration.Version, err)
		}

		// Record migration
		if _, err := s.db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", migration.Version); err != nil {
			return fmt.Errorf("failed to record migration %s: %w", migration.Version, err)
		}
	}

	return nil
}

type Migration struct {
	Version string
	SQL     string
}

func (s *Store) loadMigrations() ([]Migration, error) {
	var migrations []Migration
	
	// For now, hardcode the initial migration
	// In a full implementation, this would read from the migrations/ directory
	migrations = append(migrations, Migration{
		Version: "001_init",
		SQL: `-- Initial schema for notebook.oceanheart.ai blog engine
-- Posts table with metadata and caching
CREATE TABLE posts (
  id INTEGER PRIMARY KEY,
  slug TEXT UNIQUE NOT NULL,
  title TEXT NOT NULL,
  summary TEXT,
  html TEXT NOT NULL,
  raw_md TEXT NOT NULL,
  published_at DATETIME NOT NULL,
  updated_at DATETIME NOT NULL,
  draft BOOLEAN NOT NULL DEFAULT 0
);

-- Tags table for categorization
CREATE TABLE tags (
  id INTEGER PRIMARY KEY,
  name TEXT UNIQUE NOT NULL
);

-- Many-to-many relationship between posts and tags
CREATE TABLE post_tags (
  post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
  tag_id  INTEGER NOT NULL REFERENCES tags(id)  ON DELETE CASCADE,
  PRIMARY KEY (post_id, tag_id)
);

-- Indexes for performance
CREATE INDEX idx_posts_published ON posts(published_at DESC, draft);
CREATE INDEX idx_tags_name ON tags(name);`,
	})

	return migrations, nil
}

func (s *Store) isMigrationApplied(version string) (bool, error) {
	var count int
	err := s.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version = ?", version).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// Close closes the database connection
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Basic CRUD operations for posts
func (s *Store) GetAllPosts(includeDrafts bool) ([]Post, error) {
	query := "SELECT id, slug, title, summary, html, raw_md, published_at, updated_at, draft FROM posts"
	args := []interface{}{}
	
	if !includeDrafts {
		query += " WHERE draft = 0"
	}
	
	query += " ORDER BY published_at DESC"
	
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Slug, &p.Title, &p.Summary, &p.HTML, &p.RawMD, &p.PublishedAt, &p.UpdatedAt, &p.Draft)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}

	return posts, rows.Err()
}

func (s *Store) GetPostBySlug(slug string) (*Post, error) {
	query := "SELECT id, slug, title, summary, html, raw_md, published_at, updated_at, draft FROM posts WHERE slug = ?"
	
	var p Post
	err := s.db.QueryRow(query, slug).Scan(&p.ID, &p.Slug, &p.Title, &p.Summary, &p.HTML, &p.RawMD, &p.PublishedAt, &p.UpdatedAt, &p.Draft)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &p, nil
}

func (s *Store) UpsertPost(p *Post) error {
	query := `
		INSERT INTO posts (slug, title, summary, html, raw_md, published_at, updated_at, draft)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			title = excluded.title,
			summary = excluded.summary,
			html = excluded.html,
			raw_md = excluded.raw_md,
			published_at = excluded.published_at,
			updated_at = excluded.updated_at,
			draft = excluded.draft
	`
	
	_, err := s.db.Exec(query, p.Slug, p.Title, p.Summary, p.HTML, p.RawMD, p.PublishedAt, p.UpdatedAt, p.Draft)
	return err
}

// UpsertPosts batch upserts multiple posts
func (s *Store) UpsertPosts(posts []*Post) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Prepare statement for post upserts
	postStmt, err := tx.Prepare(`
		INSERT INTO posts (slug, title, summary, html, raw_md, published_at, updated_at, draft)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(slug) DO UPDATE SET
			title = excluded.title,
			summary = excluded.summary,
			html = excluded.html,
			raw_md = excluded.raw_md,
			published_at = excluded.published_at,
			updated_at = excluded.updated_at,
			draft = excluded.draft
	`)
	if err != nil {
		return err
	}
	defer postStmt.Close()

	for _, post := range posts {
		_, err := postStmt.Exec(post.Slug, post.Title, post.Summary, post.HTML, post.RawMD, post.PublishedAt, post.UpdatedAt, post.Draft)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetOrCreateTag gets existing tag or creates new one
func (s *Store) GetOrCreateTag(name string) (*Tag, error) {
	// Try to get existing tag
	var tag Tag
	err := s.db.QueryRow("SELECT id, name FROM tags WHERE name = ?", name).Scan(&tag.ID, &tag.Name)
	if err == nil {
		return &tag, nil
	}

	// Create new tag
	result, err := s.db.Exec("INSERT INTO tags (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &Tag{ID: int(id), Name: name}, nil
}

// LinkPostTags links a post to tags
func (s *Store) LinkPostTags(postID int, tagNames []string) error {
	if len(tagNames) == 0 {
		return nil
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Clear existing tags for this post
	_, err = tx.Exec("DELETE FROM post_tags WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	// Add new tags
	for _, tagName := range tagNames {
		tag, err := s.GetOrCreateTag(tagName)
		if err != nil {
			return err
		}

		_, err = tx.Exec("INSERT INTO post_tags (post_id, tag_id) VALUES (?, ?)", postID, tag.ID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// PopularTag represents a tag with its usage count
type PopularTag struct {
	Name  string
	Count int
}

// GetPopularTags returns the most popular tags by post count
func (s *Store) GetPopularTags(limit int) ([]PopularTag, error) {
	query := `
		SELECT t.name, COUNT(pt.post_id) as post_count
		FROM tags t
		JOIN post_tags pt ON t.id = pt.tag_id
		JOIN posts p ON pt.post_id = p.id
		WHERE p.draft = 0 
		  AND p.published_at <= datetime('now')
		GROUP BY t.name
		ORDER BY post_count DESC, t.name ASC
		LIMIT ?
	`
	
	rows, err := s.db.Query(query, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []PopularTag
	for rows.Next() {
		var tag PopularTag
		if err := rows.Scan(&tag.Name, &tag.Count); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, rows.Err()
}