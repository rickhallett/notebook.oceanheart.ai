# Change Log: Tag Navigation
## Date: 2025-09-17

## Files Modified

### internal/store/sqlite.go
- **Change**: Added PopularTag struct and GetPopularTags method
- **Rationale**: Needed to retrieve popular tags for navigation display
- **Impact**: Provides backend capability to query tags by usage count
- **Commit**: 1398a9b

### internal/http/handlers.go
- **Change**: Added getPopularTags helper method and PopularTags/ActiveTag to all handler data
- **Rationale**: Pass tag data to templates for navigation rendering
- **Impact**: All pages now have access to popular tags for navigation
- **Commit**: 31e6151

### internal/view/templates/layouts/base.html
- **Change**: Updated navigation to display popular tag links with active states
- **Rationale**: Provide visual tag navigation in header
- **Impact**: Users can now filter posts by tags from any page
- **Commit**: 387edcb

### internal/view/assets/app.css
- **Change**: Added styling for tag navigation links and responsive behavior
- **Rationale**: Visual design for tag pills and mobile optimization
- **Impact**: Tags appear as styled pills with hover/active states
- **Commit**: 0cf66d7

### internal/store/sqlite.go & internal/content/loader.go
- **Change**: Added Tags field to Post struct and tag persistence in UpsertPosts
- **Rationale**: Fix missing tag loading - tags were parsed but not saved to database
- **Impact**: Tags are now properly loaded from markdown and saved to database
- **Commit**: 9cdb4f3

## Dependencies Added/Removed
- None

## Breaking Changes
- None