# Tag Navigation PRD

**Date:** 2025-09-17  
**Feature:** Tag-based Navigation Links in Header

## Executive Summary

Add tag filtering links to the site's navigation bar to enable quick content discovery by topic. The navigation will display popular tags as clickable links that filter the index page to show only posts with the selected tag. This provides readers with intuitive topic-based browsing without requiring search or manual URL manipulation.

## Problem Statement

### Current Issues
1. **No visible tag discovery** - Tags exist in the database but are not exposed in the UI
2. **Hidden navigation paths** - The `/tag/{name}` route exists but users have no way to discover it
3. **Missed content connections** - Readers cannot easily explore related content by topic
4. **Limited navigation** - The nav bar only contains a "Home" link, underutilizing prime navigation space

### User Impact
- Readers cannot browse content by topic/category
- No way to discover what topics the blog covers
- Related content remains disconnected and harder to discover
- Poor content discoverability reduces engagement

## Requirements

### User Requirements
1. Display clickable tag links in the navigation bar
2. Show most popular/relevant tags to avoid overwhelming the navigation
3. Clicking a tag filters the homepage to show only posts with that tag
4. Visual indicator for the currently active tag filter
5. Ability to return to unfiltered view easily

### Technical Requirements
1. Query popular tags from the database (by post count)
2. Pass tag data to the base template for navigation rendering
3. Support active state styling for current tag filter
4. Maintain existing `/tag/{name}` routing pattern
5. Cache tag list appropriately for performance

### Design Requirements
1. Tags should be visually distinct but not overwhelming
2. Mobile-responsive layout (horizontal scroll or dropdown for many tags)
3. Clear visual feedback for active/hover states
4. Consistent styling with existing navigation

## Implementation Phases

### Phase 1: Core Tag Navigation
1. **Backend Tag Retrieval**
   - Add `GetPopularTags()` method to store package
   - Query top N tags by usage count
   - Return tag names and counts

2. **Template Data Structure**
   - Extend base template data with `PopularTags []Tag`
   - Include current active tag in context
   - Pass data through all handler functions

3. **Navigation Template Update**
   - Add tag links after "Home" link
   - Apply active class to current tag
   - Include "All" or similar option to clear filter

### Phase 2: Enhanced UX
1. **Tag Categorization**
   - Group special tags (cognitive-skill:*, bias:*)
   - Separate technical tags from conceptual tags
   - Consider dropdown or submenu for categories

2. **Mobile Optimization**
   - Horizontal scroll container for mobile
   - Or collapse into dropdown menu
   - Touch-friendly tap targets

3. **Visual Polish**
   - Tag pills/badges styling
   - Post count indicators (optional)
   - Smooth transitions and hover effects

## Implementation Notes

### Database Query
```sql
-- Get popular tags with post counts
SELECT t.name, COUNT(pt.post_id) as post_count
FROM tags t
JOIN post_tags pt ON t.id = pt.tag_id
JOIN posts p ON pt.post_id = p.id
WHERE p.draft = 0 
  AND p.published_at <= datetime('now')
GROUP BY t.name
ORDER BY post_count DESC
LIMIT 10;
```

### Template Structure
```html
<!-- Update base.html navigation -->
<nav class="navigation">
    <a href="/" class="{{if not .ActiveTag}}active{{end}}">Home</a>
    {{range .PopularTags}}
        <a href="/tag/{{.Name}}" 
           class="tag-link {{if eq $.ActiveTag .Name}}active{{end}}">
           {{.Name}}
        </a>
    {{end}}
</nav>
```

### Handler Enhancement
```go
// Add to Server struct or context
type NavContext struct {
    PopularTags []Tag
    ActiveTag   string
}

// Load tags once at startup or cache with TTL
func (s *Server) loadPopularTags() []Tag {
    return s.store.GetPopularTags(10) // Top 10 tags
}
```

### CSS Styling Example
```css
.navigation {
    display: flex;
    gap: 1rem;
    flex-wrap: wrap;
    align-items: center;
}

.tag-link {
    padding: 0.25rem 0.75rem;
    border-radius: 1rem;
    background: var(--tag-bg);
    font-size: 0.9em;
}

.tag-link.active,
.tag-link:hover {
    background: var(--tag-active-bg);
    color: var(--tag-active-color);
}

/* Mobile */
@media (max-width: 768px) {
    .navigation {
        overflow-x: auto;
        flex-wrap: nowrap;
    }
}
```

## Security Considerations

1. **Input Validation**
   - Sanitize tag names in URLs to prevent XSS
   - Validate tag exists before querying posts
   - URL decode tag names properly for special characters

2. **Performance**
   - Cache popular tags list (refresh hourly or on content change)
   - Limit number of tags shown to prevent DOM bloat
   - Consider lazy loading for tag dropdown if implemented

## Success Metrics

1. **Engagement Metrics**
   - Click-through rate on tag navigation links
   - Average posts viewed per session after implementation
   - Bounce rate reduction on tagged content pages

2. **Technical Metrics**
   - Page load time remains under 200ms
   - Navigation render completes in single paint
   - No layout shift from tag loading

## Future Enhancements

1. **Dynamic Tag Cloud**
   - Variable sizing based on popularity
   - Color coding for tag categories
   - Animated transitions on hover

2. **Tag Search/Filter**
   - Type-ahead search in tag list
   - Multi-tag filtering (AND/OR operations)
   - Tag exclusion filters

3. **Personalization**
   - Remember user's preferred tags
   - Suggest tags based on reading history
   - Custom tag subscriptions via RSS

4. **Tag Management**
   - Admin UI for tag aliases/synonyms
   - Tag merging and cleanup tools
   - Tag usage analytics dashboard

## Notes

- Start with a simple horizontal list of top 5-10 tags
- The existing `/tag/{name}` route handler can be reused without modification
- Consider special handling for cognitive-skill:* tags which already have distinct visual treatment
- Ensure graceful degradation if no tags exist in the system
- Keep the implementation minimal initially - avoid over-engineering the tag system