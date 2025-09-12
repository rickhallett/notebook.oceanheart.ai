package content

import (
	"os"
	"testing"

	"notebook.oceanheart.ai/internal/store"
)

func TestContentPipelineIntegration(t *testing.T) {
	// Create temporary database
	tempDB := "test_integration.db"
	defer os.Remove(tempDB)

	db := store.MustOpen(tempDB)
	defer db.Close()

	// Test content loading
	loader := NewLoader("../../content")
	posts, err := loader.LoadAll()
	if err != nil {
		t.Fatalf("LoadAll failed: %v", err)
	}

	if len(posts) == 0 {
		t.Skip("No content files found, skipping integration test")
	}

	// Test database storage
	err = db.UpsertPosts(posts)
	if err != nil {
		t.Fatalf("UpsertPosts failed: %v", err)
	}

	// Verify posts were stored
	storedPosts, err := db.GetAllPosts(true) // Include drafts for testing
	if err != nil {
		t.Fatalf("GetAllPosts failed: %v", err)
	}

	if len(storedPosts) != len(posts) {
		t.Errorf("Expected %d posts in database, got %d", len(posts), len(storedPosts))
	}

	// Test retrieving a specific post
	if len(storedPosts) > 0 {
		post := storedPosts[0]
		retrieved, err := db.GetPostBySlug(post.Slug)
		if err != nil {
			t.Fatalf("GetPostBySlug failed: %v", err)
		}

		if retrieved == nil {
			t.Error("Expected to find post by slug")
		} else if retrieved.Title != post.Title {
			t.Errorf("Retrieved post title mismatch: got %s, expected %s", retrieved.Title, post.Title)
		}
	}

	// Verify HTML was generated
	for _, post := range storedPosts {
		if len(post.HTML) == 0 {
			t.Errorf("Post %s has no HTML content", post.Slug)
		}
		if len(post.RawMD) == 0 {
			t.Errorf("Post %s has no raw markdown content", post.Slug)
		}
	}
}