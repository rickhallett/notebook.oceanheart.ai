package feed

import (
	"encoding/xml"
	"fmt"
	"time"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/store"
)

// AtomFeed represents an Atom 1.0 feed structure
type AtomFeed struct {
	XMLName xml.Name   `xml:"feed"`
	Xmlns   string     `xml:"xmlns,attr"`
	Title   string     `xml:"title"`
	Link    []AtomLink `xml:"link"`
	ID      string     `xml:"id"`
	Updated string     `xml:"updated"`
	Author  AtomAuthor `xml:"author"`
	Entry   []AtomEntry `xml:"entry"`
}

// AtomLink represents a link in an Atom feed
type AtomLink struct {
	Href string `xml:"href,attr"`
	Rel  string `xml:"rel,attr,omitempty"`
	Type string `xml:"type,attr,omitempty"`
}

// AtomAuthor represents the author of an Atom feed
type AtomAuthor struct {
	Name string `xml:"name"`
	URI  string `xml:"uri,omitempty"`
}

// AtomEntry represents an entry in an Atom feed
type AtomEntry struct {
	Title   string     `xml:"title"`
	Link    []AtomLink `xml:"link"`
	ID      string     `xml:"id"`
	Updated string     `xml:"updated"`
	Summary string     `xml:"summary,omitempty"`
	Content AtomContent `xml:"content"`
}

// AtomContent represents the content of an Atom entry
type AtomContent struct {
	Type string `xml:"type,attr"`
	Body string `xml:",chardata"`
}

// GenerateAtom creates an Atom 1.0 feed from posts
func GenerateAtom(posts []*store.Post, cfg *config.Config) ([]byte, error) {
	if len(posts) == 0 {
		return nil, fmt.Errorf("no posts available for feed")
	}

	// Limit to latest 20 posts
	feedPosts := posts
	if len(posts) > 20 {
		feedPosts = posts[:20]
	}

	// Create feed structure
	feed := AtomFeed{
		Xmlns: "http://www.w3.org/2005/Atom",
		Title: cfg.SiteTitle,
		Link: []AtomLink{
			{Href: cfg.SiteBaseURL, Rel: "alternate", Type: "text/html"},
			{Href: cfg.SiteBaseURL + "/feed.xml", Rel: "self", Type: "application/atom+xml"},
		},
		ID:      cfg.SiteBaseURL + "/",
		Updated: formatAtomDate(feedPosts[0].UpdatedAt),
		Author: AtomAuthor{
			Name: "Oceanheart",
			URI:  cfg.SiteBaseURL,
		},
	}

	// Add entries
	for _, post := range feedPosts {
		// Skip drafts in production
		if post.Draft && cfg.Environment != "dev" {
			continue
		}

		entry := AtomEntry{
			Title: post.Title,
			Link: []AtomLink{
				{Href: cfg.SiteBaseURL + "/p/" + post.Slug, Rel: "alternate", Type: "text/html"},
			},
			ID:      cfg.SiteBaseURL + "/p/" + post.Slug,
			Updated: formatAtomDate(post.UpdatedAt),
			Summary: post.Summary,
			Content: AtomContent{
				Type: "html",
				Body: post.HTML,
			},
		}
		feed.Entry = append(feed.Entry, entry)
	}

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(feed, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal atom feed: %w", err)
	}

	// Add XML declaration
	xmlWithDeclaration := []byte(xml.Header + string(xmlData))
	return xmlWithDeclaration, nil
}

// formatAtomDate formats time string to RFC3339 for Atom feeds
func formatAtomDate(dateStr string) string {
	// Parse the date string (should already be in RFC3339 format)
	t, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Fallback to current time if parsing fails
		t = time.Now()
	}
	return t.Format(time.RFC3339)
}