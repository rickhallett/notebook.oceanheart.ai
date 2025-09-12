package store

import (
	"os"
	"testing"
)

func TestStoreMigrations(t *testing.T) {
	// Create temporary database file
	tempDB := "test_notebook.db"
	defer os.Remove(tempDB)

	// Test database creation and migration
	store := MustOpen(tempDB)
	defer store.Close()

	// Verify database connection works
	if err := store.db.Ping(); err != nil {
		t.Fatalf("Database ping failed: %v", err)
	}

	// Verify tables exist
	tables := []string{"posts", "tags", "post_tags", "schema_migrations"}
	for _, table := range tables {
		var count int
		query := "SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?"
		err := store.db.QueryRow(query, table).Scan(&count)
		if err != nil {
			t.Fatalf("Failed to check table %s: %v", table, err)
		}
		if count != 1 {
			t.Errorf("Table %s not found", table)
		}
	}

	// Verify migration was recorded
	var migrationCount int
	err := store.db.QueryRow("SELECT COUNT(*) FROM schema_migrations WHERE version='001_init'").Scan(&migrationCount)
	if err != nil {
		t.Fatalf("Failed to check migration record: %v", err)
	}
	if migrationCount != 1 {
		t.Error("Migration 001_init was not recorded")
	}
}

func TestBasicPostOperations(t *testing.T) {
	// Create temporary database
	tempDB := "test_posts.db"
	defer os.Remove(tempDB)

	store := MustOpen(tempDB)
	defer store.Close()

	// Test getting posts from empty database
	posts, err := store.GetAllPosts(false)
	if err != nil {
		t.Fatalf("GetAllPosts failed: %v", err)
	}
	if len(posts) != 0 {
		t.Errorf("Expected 0 posts, got %d", len(posts))
	}

	// Test getting post by slug from empty database
	post, err := store.GetPostBySlug("nonexistent")
	if err != nil {
		t.Fatalf("GetPostBySlug failed: %v", err)
	}
	if post != nil {
		t.Error("Expected nil for nonexistent post")
	}

	// Test upserting a post
	testPost := &Post{
		Slug:        "test-post",
		Title:       "Test Post",
		Summary:     "A test post",
		HTML:        "<p>Test content</p>",
		RawMD:       "Test content",
		PublishedAt: "2025-09-12T10:00:00Z",
		UpdatedAt:   "2025-09-12T10:00:00Z",
		Draft:       false,
	}

	err = store.UpsertPost(testPost)
	if err != nil {
		t.Fatalf("UpsertPost failed: %v", err)
	}

	// Test retrieving the post
	retrieved, err := store.GetPostBySlug("test-post")
	if err != nil {
		t.Fatalf("GetPostBySlug failed: %v", err)
	}
	if retrieved == nil {
		t.Fatal("Expected post to be found")
	}
	if retrieved.Title != "Test Post" {
		t.Errorf("Expected title 'Test Post', got %s", retrieved.Title)
	}

	// Test getting all posts includes our post
	posts, err = store.GetAllPosts(false)
	if err != nil {
		t.Fatalf("GetAllPosts failed: %v", err)
	}
	if len(posts) != 1 {
		t.Errorf("Expected 1 post, got %d", len(posts))
	}
}