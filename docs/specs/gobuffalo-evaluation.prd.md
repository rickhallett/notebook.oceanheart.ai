# GoBuffalo Migration DX Evaluation (PRD)

## Executive Summary
Evaluate replacing the current stdlib-based Go app with GoBuffalo to improve developer experience (DX). Consider routing, templating, assets, database, testing, publishing workflow, and deployment. Given today’s scope (Markdown blog + SQLite cache + feeds/sitemap), Buffalo provides limited immediate DX gain, but becomes valuable if we plan an admin UI, auth, richer views, or complex assets.

## Problem Statement
We need to assess whether adopting Buffalo will materially reduce friction for: local development (live reload, scaffolds), adding features (admin, forms, auth), and operating the app (migrations, assets), without harming deploy simplicity (single binary, Caddy, systemd) and the content publishing workflow.

## Current Feature Set (Baseline)
- Content: Markdown + YAML front matter; SQLite caching; admin-less.
- HTTP: stdlib mux; inline templates; middleware chain; SEO meta; feed/sitemap.
- Ops: single binary; systemd; Caddy reverse proxy; rsync publish; optional reload endpoint.

## Requirements for Comparison
- DX: generators, hot reload, routing ergonomics, templates, asset pipeline.
- Data: migrations, schema changes, test DB management.
- Publishing: write locally → preview → publish remote with minimal steps.
- Deploy: continue systemd + Caddy; keep small operational footprint.

## Analysis: GoBuffalo vs Current
- Routing & Controllers
  - Buffalo: resource routing, path helpers, middleware stacking → faster CRUD/admin.
  - Current: simple mux adequate for blog; minimal ceremony.
- Templates & Assets
  - Buffalo: Plush templates, asset pipeline (fingerprinting, minify), dev live reload.
  - Current: inline templates; no pipeline. File-based templates + a light watcher would close most DX gaps without a framework.
- Database & Migrations
  - Buffalo Pop: schema management, migration CLI, fixtures → better DX for feature work.
  - Current: hand-rolled migrations (single init); fine for read-mostly blog.
- Testing
  - Buffalo: helpers for actions/models; good for web apps.
  - Current: std `testing`, already green and fast.
- Publishing Workflow
  - Buffalo: `buffalo dev` for live reload; `buffalo build` for assets + binary. Publish still rsync/systemd.
  - Current: go build + rsync already simple; no asset complexity today.
- Deployment
  - Buffalo adds dependencies, larger binary, assets to serve. Still compatible with Caddy + systemd.
- Security & Sessions
  - Buffalo includes CSRF/sessions scaffolding useful for future admin.

## Recommendation
- Short term (blog-only): stay on stdlib. Add file-based templates and a live-reload tool (`air`/`reflex`) to boost DX with minimal complexity.
- If adding admin UI/auth/forms soon: consider migrating to Buffalo to leverage routing, scaffolds, Pop, CSRF, and asset pipeline.

## Success Metrics
- Time to add a basic admin “new post” UI: Current: medium-high; Buffalo: low-medium.
- Onboarding time for new devs: Current: low; Buffalo: low-medium (framework learning curve).
- Build/deploy time: Current: minimal; Buffalo: slightly higher (assets + larger binary).

## Migration Approach (if chosen)
- Phase 1: Create Buffalo app skeleton; mount existing content rendering + store package as internal modules.
- Phase 2: Replace routing/views gradually; move templates to `templates/`, static to `public/`.
- Phase 3: Introduce Pop migrations; map existing schema; keep SQLite.
- Phase 4: Add admin screens with CSRF/session; preserve `/feed.xml` and `/sitemap.xml` routes.

## Risks & Mitigations
- Increased complexity/deps → Start with hybrid: adopt only where DX gains are clear (admin).
- Performance/binary size → Validate with benchmarks; still acceptable for target.
- Team familiarity → Document patterns; keep changes incremental.

## Decision
- Default: remain on current stack; re-evaluate when admin UI is prioritized.
