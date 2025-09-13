# File-Based Templates (PRD)

## Executive Summary
Replace inline templates with file-based HTML templates to improve DX, maintainability, and theming. Keep the current data model and handlers; introduce a template loader with dev-time reload and production caching. No change to public routes.

## Problem Statement
Inline templates in `internal/http/handlers.go` hinder readability, reuse (partials), and styling. We need conventional layouts/partials, easier iteration, and the ability to theme without code changes.

## Requirements
- Use Go `html/template` with a clear directory structure.
- Dev: auto-reload templates on each request; Prod: compile once at startup.
- Layout/partials support; page templates for home, post, tag.
- Safe rendering of post HTML (`template.HTML` for trusted renderer output).
- Keep route outputs and meta tags compatible with current HTML.
- Simple asset pathing for CSS (served from `/static/`).

## Implementation Phases
- Phase 1: Template structure and loader
  - `internal/view/templates/`
    - `layouts/base.html`
    - `partials/{head.html,nav.html,footer.html}`
    - `pages/{home.html,post.html,tag.html}`
  - Loader: `LoadTemplates(glob ...string)`; inject funcs (date formatting, safeHTML).
- Phase 2: Handlers -> ExecuteTemplate
  - Replace inline template usage with `ExecuteTemplate(w, "pages/home.html", data)` etc.
  - Preserve meta, SEO, and link tags.
- Phase 3: Dev reload & prod cache
  - If `cfg.IsDev()`: parse on each request; else parse at startup.
- Phase 4: Static assets
  - Add `internal/view/assets/` (CSS), serve at `/static/`.

## Implementation Notes
- Data shape: introduce `PageData` with fields already passed today (`Title`, `Description`, `CanonicalURL`, `BaseURL`, `SiteTitle`, `IsPost`, `PublishedAt`, `UpdatedAt`, `Content`, `Posts`).
- Template funcs: `safeHTML(string) template.HTML`, `formatDateRFC3339(string)`.
- Parsing: `template.ParseFS` or `template.ParseGlob("internal/view/templates/**/*.html")`.
- Error handling: on parse/exec error, log and return `500` with a simple fallback.
- Tests: golden/snapshot tests using `httptest` to verify head/meta, body structure, and key content.

## Security Considerations
- Only mark renderer-produced post content as safe. Escape all other fields by default.
- Keep security headers middleware unchanged.

## Success Metrics
- Handlers reduced in size; template changes do not require Go recompile in dev.
- Add a new partial or tweak head tags without touching Go code.
- All existing tests pass; add coverage for template rendering paths.

## Future Enhancements
- Theme packs under `themes/` with overridable templates.
- Live reload with fsnotify.
- Layout blocks for per-page CSS/JS.
