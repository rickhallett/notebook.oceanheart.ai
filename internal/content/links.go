package content

import (
	"regexp"
	"strings"
)

// ProcessExternalLinks adds security attributes to external links in HTML
func ProcessExternalLinks(html string, baseURL string) string {
	// Regex to match <a href="..."> tags
	linkRegex := regexp.MustCompile(`<a\s+([^>]*\s+)?href="([^"]+)"([^>]*)>`)
	
	return linkRegex.ReplaceAllStringFunc(html, func(match string) string {
		// Extract href URL from the match
		hrefRegex := regexp.MustCompile(`href="([^"]+)"`)
		hrefMatch := hrefRegex.FindStringSubmatch(match)
		if len(hrefMatch) < 2 {
			return match // Return original if we can't extract URL
		}
		
		url := hrefMatch[1]
		
		// Check if it's an external link
		if isExternalLink(url, baseURL) {
			// Check if it already has target and rel attributes
			hasTarget := strings.Contains(match, "target=")
			hasRel := strings.Contains(match, "rel=")
			
			// Add attributes if missing
			var additions []string
			if !hasTarget {
				additions = append(additions, `target="_blank"`)
			}
			if !hasRel {
				additions = append(additions, `rel="noopener noreferrer"`)
			}
			
			if len(additions) > 0 {
				// Insert additions before the closing >
				closingIndex := strings.LastIndex(match, ">")
				if closingIndex != -1 {
					return match[:closingIndex] + " " + strings.Join(additions, " ") + match[closingIndex:]
				}
			}
		}
		
		return match
	})
}

// isExternalLink determines if a URL is external to the base domain
func isExternalLink(url, baseURL string) bool {
	// Skip relative URLs, anchors, and protocol-relative URLs
	if strings.HasPrefix(url, "/") || 
	   strings.HasPrefix(url, "#") || 
	   strings.HasPrefix(url, "//") {
		return false
	}
	
	// Skip mailto and other non-HTTP protocols
	if strings.HasPrefix(url, "mailto:") || 
	   strings.HasPrefix(url, "tel:") ||
	   strings.HasPrefix(url, "javascript:") {
		return false
	}
	
	// If URL doesn't start with http/https, assume it's relative
	if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
		return false
	}
	
	// Extract domain from base URL
	baseURLLower := strings.ToLower(baseURL)
	urlLower := strings.ToLower(url)
	
	// Remove protocol from base URL to get domain
	baseDomain := baseURLLower
	if strings.HasPrefix(baseDomain, "https://") {
		baseDomain = strings.TrimPrefix(baseDomain, "https://")
	} else if strings.HasPrefix(baseDomain, "http://") {
		baseDomain = strings.TrimPrefix(baseDomain, "http://")
	}
	
	// Remove trailing path from base domain
	if idx := strings.Index(baseDomain, "/"); idx != -1 {
		baseDomain = baseDomain[:idx]
	}
	
	// Check if the URL starts with our base domain
	return !strings.HasPrefix(urlLower, "https://"+baseDomain) && 
		   !strings.HasPrefix(urlLower, "http://"+baseDomain)
}