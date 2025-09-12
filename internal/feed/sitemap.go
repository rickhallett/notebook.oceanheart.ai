package feed

import (
	"encoding/xml"
	"fmt"
	"time"

	"notebook.oceanheart.ai/internal/config"
	"notebook.oceanheart.ai/internal/store"
)

// SitemapIndex represents a sitemap.xml structure
type SitemapIndex struct {
	XMLName xml.Name `xml:"urlset"`
	Xmlns   string   `xml:"xmlns,attr"`
	URLs    []SitemapURL `xml:"url"`
}

// SitemapURL represents a URL in a sitemap
type SitemapURL struct {
	Loc        string `xml:"loc"`
	LastMod    string `xml:"lastmod,omitempty"`
	ChangeFreq string `xml:"changefreq,omitempty"`
	Priority   string `xml:"priority,omitempty"`
}

// GenerateSitemap creates an XML sitemap from posts and static pages
func GenerateSitemap(posts []*store.Post, cfg *config.Config) ([]byte, error) {
	sitemap := SitemapIndex{
		Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
	}

	// Add home page
	sitemap.URLs = append(sitemap.URLs, SitemapURL{
		Loc:        cfg.SiteBaseURL + "/",
		LastMod:    formatSitemapDate(time.Now()),
		ChangeFreq: "daily",
		Priority:   "1.0",
	})

	// Add posts
	for _, post := range posts {
		// Skip drafts in production
		if post.Draft && cfg.Environment != "dev" {
			continue
		}

		// Parse updated time for lastmod
		lastMod := formatSitemapDate(time.Now())
		if t, err := time.Parse(time.RFC3339, post.UpdatedAt); err == nil {
			lastMod = formatSitemapDate(t)
		}

		sitemap.URLs = append(sitemap.URLs, SitemapURL{
			Loc:        cfg.SiteBaseURL + "/p/" + post.Slug,
			LastMod:    lastMod,
			ChangeFreq: "monthly",
			Priority:   "0.8",
		})
	}

	// Add feed URL
	sitemap.URLs = append(sitemap.URLs, SitemapURL{
		Loc:        cfg.SiteBaseURL + "/feed.xml",
		LastMod:    formatSitemapDate(time.Now()),
		ChangeFreq: "daily",
		Priority:   "0.5",
	})

	// Marshal to XML
	xmlData, err := xml.MarshalIndent(sitemap, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal sitemap: %w", err)
	}

	// Add XML declaration
	xmlWithDeclaration := []byte(xml.Header + string(xmlData))
	return xmlWithDeclaration, nil
}

// formatSitemapDate formats time to W3C datetime format for sitemaps
func formatSitemapDate(t time.Time) string {
	return t.Format("2006-01-02T15:04:05-07:00")
}