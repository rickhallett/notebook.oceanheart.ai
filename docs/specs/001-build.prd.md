# notebook.oceanheart.ai - Development Breakdown PRD

## Project Analysis

**Project**: notebook.oceanheart.ai - A minimalist Go + HTMX + SQLite blog engine for "learning in public"

**Core Architecture**: 
- Backend: Go (net/http) 
- Frontend: HTMX with progressive enhancement
- Database: SQLite with auto-migrations
- Content: Markdown-first with front matter
- Deployment: Single binary + Docker

**Key Features**:
- Markdown posts with front matter
- HTMX-powered search/filter and pagination
- Tags & archives, RSS/Atom feeds
- Drafts mode for development
- Code highlighting with Chroma
- "Psychology twist" metadata with cognitive skill/bias tags
- Admin endpoints for reindexing and cache management

## Development Phase Breakdown

### Phase 00: Core Infrastructure & Database Foundation
**Priority**: Critical Foundation
**Complexity**: Medium
**Estimated Effort**: 3-4 days

**Scope**:
- SQLite schema setup with migrations
- Basic project structure and Go modules
- Database connection and migration system
- Environment configuration management
- Health check endpoint

**Dependencies**: None

**Testing Strategy**:
- Database connection and migration tests
- Environment variable loading tests
- Basic health check validation

**Success Criteria**:
- Database creates and migrates successfully
- Environment variables load correctly
- Health check returns 200 status
- Project structure follows specified layout

**Key Deliverables**:
- `/migrations/001_init.sql` with posts, tags, and post_tags tables
- `/internal/store/sqlite.go` with basic CRUD operations
- `/cmd/notebook/main.go` with basic server setup
- Environment configuration system

---

### Phase 01: Content Loading & Markdown Processing
**Priority**: Critical Foundation  
**Complexity**: Medium-High
**Estimated Effort**: 4-5 days

**Scope**:
- Markdown file parsing with front matter
- Content loader that scans `/content` directory
- Markdown to HTML rendering with goldmark
- Code highlighting integration with Chroma
- Basic caching mechanism

**Dependencies**: Phase 00 (database foundation)

**Testing Strategy**:
- Markdown parsing with various front matter formats
- HTML rendering output validation
- Code highlighting for multiple languages
- Content loading from filesystem

**Success Criteria**:
- Parses markdown files with YAML front matter
- Renders markdown to HTML with code highlighting
- Loads all content from `/content` directory
- Caches rendered HTML in database

**Key Deliverables**:
- `/internal/content/loader.go` 
- `/internal/content/render.go`
- Content parsing and rendering pipeline
- Basic template system foundation

---

### Phase 02: HTTP Handlers & Basic Routing
**Priority**: Critical Core
**Complexity**: Medium
**Estimated Effort**: 3-4 days

**Scope**:
- HTTP server setup with net/http
- Basic route handlers (home, post, tag)
- Template rendering system
- Middleware for logging and gzip
- Static asset serving

**Dependencies**: Phase 01 (content processing)

**Testing Strategy**:
- Route response validation
- Template rendering tests
- Middleware functionality
- Static asset delivery

**Success Criteria**:
- All basic routes return proper responses
- Templates render with dynamic content
- Middleware processes requests correctly
- Static assets serve properly

**Key Deliverables**:
- `/internal/http/handlers.go`
- `/internal/http/middleware.go`
- `/internal/view/templates/` base templates
- Basic CSS framework

---

### Phase 03: HTMX Integration & Search
**Priority**: High Value Feature
**Complexity**: Medium-High
**Estimated Effort**: 4-5 days

**Scope**:
- HTMX search functionality with partial results
- Search handler returning template fragments
- Basic pagination with "Load more" functionality
- HTMX-specific template partials

**Dependencies**: Phase 02 (basic routing)

**Testing Strategy**:
- Search query processing and results
- HTMX partial template rendering
- Pagination logic validation
- Client-side HTMX behavior testing

**Success Criteria**:
- Search returns filtered results via HTMX
- Pagination loads additional content dynamically
- HTMX templates render correctly in browser
- Search is responsive and fast

**Key Deliverables**:
- Search handler with partial template support
- HTMX pagination implementation
- `/internal/view/templates/partials/` directory
- JavaScript-free progressive enhancement

---

### Phase 04: Feeds & SEO Features
**Priority**: Medium Value Feature
**Complexity**: Low-Medium
**Estimated Effort**: 2-3 days

**Scope**:
- RSS/Atom feed generation
- Sitemap.xml generation
- Meta tags for SEO (title, description, published time)
- External link handling with proper attributes

**Dependencies**: Phase 02 (basic routing)

**Testing Strategy**:
- Feed XML validation and structure
- Sitemap XML compliance
- Meta tag presence and accuracy
- External link attribute verification

**Success Criteria**:
- Valid RSS/Atom feed with latest 20 posts
- Compliant sitemap.xml with all post URLs
- Proper SEO meta tags on all pages
- External links open safely in new tabs

**Key Deliverables**:
- Feed generation handler
- Sitemap generation functionality
- SEO meta tag templates
- Link processing enhancement

---

### Phase 05: Admin Interface & Content Management
**Priority**: Medium Operational Feature
**Complexity**: Low-Medium  
**Estimated Effort**: 2-3 days

**Scope**:
- Admin reindex endpoint
- Cache flush functionality
- Environment-based admin access control
- Content reload without restart

**Dependencies**: Phase 01 (content loading)

**Testing Strategy**:
- Reindex functionality validation
- Cache clearing verification  
- Environment-based access control
- Content hot-reload testing

**Success Criteria**:
- Reindex updates content without restart
- Cache flush clears stored HTML
- Admin endpoints restricted in production
- Content changes reflect immediately

**Key Deliverables**:
- `/admin/reindex` and `/admin/flush-cache` endpoints
- Environment-based security controls
- Hot content reload functionality

---

### Phase 06: Psychology Tags & Cognitive Ribbons
**Priority**: Low-Medium Unique Feature
**Complexity**: Low
**Estimated Effort**: 2-3 days

**Scope**:
- Special handling for `cognitive-skill:*` and `bias:*` tags
- Colored ribbon rendering for psychology tags
- Optional legend page for cognitive ribbons
- Enhanced tag styling and categorization

**Dependencies**: Phase 02 (basic routing), Phase 03 (search for tag filtering)

**Testing Strategy**:
- Psychology tag detection and parsing
- Ribbon visual rendering validation
- Tag categorization and filtering
- Legend page functionality

**Success Criteria**:
- Psychology tags render as colored ribbons
- Regular tags display normally
- Tag filtering works for all tag types
- Legend page explains psychology tags

**Key Deliverables**:
- Psychology tag detection logic
- CSS ribbon styling system
- Tag categorization templates
- Optional cognitive ribbons legend page

---

### Phase 07: Production & Deployment Features
**Priority**: High Operational
**Complexity**: Medium
**Estimated Effort**: 3-4 days

**Scope**:
- Production build process and binary optimization
- Docker containerization
- Rate limiting for search endpoints
- Enhanced logging and monitoring
- Cache headers and ETag support

**Dependencies**: All previous phases (complete system)

**Testing Strategy**:
- Production build verification
- Docker container functionality
- Rate limiting behavior
- Cache header validation
- Logging format and completeness

**Success Criteria**:
- Single binary builds and runs in production
- Docker container works with volume mounting
- Rate limiting protects against abuse
- Proper cache headers improve performance
- Structured logging provides operational insight

**Key Deliverables**:
- Makefile with build targets
- Dockerfile with multi-stage build
- Rate limiting middleware
- Enhanced logging and monitoring
- Cache optimization features

---

## Implementation Strategy

**Development Principles**:
- Test-driven development for core functionality
- Defensive programming with proper error handling
- Simplicity over elegance - avoid over-engineering
- Progressive enhancement - ensure basic functionality without JavaScript
- Strict adherence to specified requirements only

**Testing Checkpoints**:
Each phase includes human-testable milestones where functionality can be verified through:
- Browser testing for user-facing features
- Manual content creation and management
- Performance validation with sample data
- Admin interface functionality verification

**Key Dependencies**:
- Phases 00-01 are foundational and must be completed first
- Phase 02 enables all user-facing functionality  
- Phases 03-06 can be developed in parallel after Phase 02
- Phase 07 requires the complete system for deployment optimization

This breakdown creates natural development boundaries where each phase delivers meaningful value and can be thoroughly tested before proceeding to the next phase. The structure prioritizes getting a working blog engine operational quickly while leaving enhancement features for later phases.