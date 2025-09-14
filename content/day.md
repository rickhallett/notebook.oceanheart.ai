commit b6f60abc90c2340419fb62c459ca020dcb9fa84e
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:13:59 2025 +0100

    docs(prd): add Phase 5 (structured logging), Phase 6 (metrics exporter), Phase 7 (error reporting) PRDs

commit 31b47c44dbc07f304e74c16b0206d7263b27426e
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:09:58 2025 +0100

    docs(prd): add presets for Phase 5 (logging), Phase 6 (metrics), Phase 7 (error reporting) to create-prd.md

commit 43de7b5883038bbe38df58aa3ae1f3e3e49d05df
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:55:38 2025 +0100

    feat(admin): add basic admin interface with RBAC, user list, role toggle, delete; minimal admin layout; docs and change log for Phase 4

commit b9fae72e4594c0572780ea172395e3b6af9ed72c
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:38:02 2025 +0100

    feat(auth): add JwtAuthentication concern with cookie helpers and helpers for API & web flows

commit d69ce8307c779abc1137d14392ef4460811cf64b
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:27:46 2025 +0100

    feat(auth): add JWT cookie helpers, API endpoints and CORS; integrate with SessionsController; sign-out redirects to home

commit 0537e83200733281a07b36b593fe5f65ebc135a3
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:27:40 2025 +0100

    feat(ui): switch to monochrome terminal theme; remove duplicate banners; add spacing/centering; update titles/prompts

commit 168a4f776bbdafd2599f84f327df180374b50ceb
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:27:33 2025 +0100

    docs(prd): add auth-ui-polish PRD, report, change log; update create-prd and implement-prd presets for terminal auth UI

commit de26f794ae9184af6ae67327cc68c58490a30587
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:27:33 2025 +0100

    docs: add AGENTS.md contributor guide (Repository Guidelines)

---

commit 890adea223fa235dff1d2222e566476fbc0b6e10
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 15:02:20 2025 +0100

    CHORE: Update package.json formatting
    
    - Minify package.json formatting for consistency
    - Add @flydotio/dockerfile as dev dependency

commit 7a2eae97661ea2ac9d28afdc14bb785d83bc33bb
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 15:02:11 2025 +0100

    CHORE: Add deployment infrastructure for Fly.io
    
    - Add multi-stage Dockerfile for Go blog engine with CA certificates
    - Add fly.toml configuration for London region deployment
    - Add .dockerignore for optimized build context
    - Add GitHub Actions workflow for automated deployment
    - Configure for port 8080 with auto-scaling machines

commit b4f63edcc1c8358de4a4fe2e0ea7810a0200e300
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 15:01:29 2025 +0100

    FEAT: Add Turso database support with fallback to SQLite
    
    - Add libsql-client-go dependency for Turso integration
    - Update store layer to support both Turso (production) and SQLite (local dev)
    - Use environment variables DB_URL and DB_AUTH_TOKEN for Turso connection
    - Maintain backward compatibility with existing SQLite setup

commit fe83892228f77bc5af9e1f8eec1d94ef1d7d6c32
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 12:12:12 2025 +0100

    docs(ops): add production deploy strategy on Fly.io and Turso migration plan with libSQL examples

commit 6a3930e0e9224acacded6be55f59da39a282bd4f
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:54:31 2025 +0100

    feat(dev): add hot-reload dev script and disable caching for static assets in dev
    
    - scripts/dev.sh supports watchexec/reflex/entr
    - Auto content reload on markdown changes via /admin/reload
    - Server restart on Go code changes
    - Set Cache-Control: no-store for /static assets in dev

commit 2d97ec0fc297b9eadb313685faa5d87dc2b5cdc2
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:45:31 2025 +0100

    docs: add blog workflow guide and PRDs for templates and GoBuffalo evaluation

commit 57afb6f34d9a17561e89872aa62a4f0ce9dd8ecd
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:45:31 2025 +0100

    docs(ops): add deployment strategy, Caddy and systemd samples, and publish script
    
    - docs/specs/deployment-strategy.prd.md with systemd, Caddy, rsync workflow
    - docs/ops/Caddyfile and docs/ops/notebook.service examples
    - scripts/publish.sh to rsync content and binary then restart service
    - caddy-setup.md for local reverse proxy on <project-name>.lvh.me

commit 186e7d2d067747cce1bb6bbd9eae8f4ec75c23a2
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:45:31 2025 +0100

    feat(app): switch to file-based templates, static assets, and admin reload endpoint
    
    - Add template manager and file-based templates under internal/view
    - Refactor handlers to render via base + page content fragments
    - Serve static assets from /static (app.css)
    - Add /admin/reload endpoint and token-based auth (RELOAD_TOKEN)
    - Wire reload route and extend config (ReloadToken, AllowReload)

commit d6becb95674f049ab8ff86b5efd1219e06ca8fe2
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:41:49 2025 +0100

    chore(git): ignore theme(s) and untrack theme subdirectories

---

commit 4c84dee1258b2b79c7a355109e0ae94b965fa98b
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 17:01:01 2025 +0100

    REFACTOR: Improve Docker build and development environment
    
    This commit introduces several improvements to the Docker setup, including:
    
    - Optimizing the Dockerfile for better caching and smaller image size.
    - Adding development dependencies and pre-collecting static files in the build.
    - Refining the docker-compose.yml for a better development experience, including a Caddy service.
    - Making the entrypoint script more robust.

commit e1e8fb35c2ccff45e561c58e95afadff6db07556
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:48:24 2025 +0100

    FIX: Ensure container uses venv Python; run manage.py and gunicorn from venv

commit 56ba11b430db20d7d2f67cead59908be462f31f5
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:37:24 2025 +0100

    BUILD: Fix Docker frontend stage to build only frontend (avoid backend collectstatic in builder)

commit e22dd623c2d3979b4a45bc1fb7af6d235b790cf6
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:49:43 2025 +0100

    CHORE: Switch dev ports to 8881 (frontend) and 8888 (Django); update Caddy, scripts, CORS, compose, and docs

commit 9375b6871a71adab912ff972644f3225b2390c01
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:21:28 2025 +0100

    DOCS: Add contributor guide and deployment PRD; CHORE: Add Caddy configs for local/prod

---

commit d08f91e9610a9206c62351ca9193c7804ffcb703
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 20:38:39 2025 +0100

    poc: app layout draft

commit 4b90d36fa35f403acf2532fe66832baef0d94519
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:39:17 2025 +0100

    DOCS: add decision log template and seed ADRs for Phases 2–3 (DSL, validation, autosave, LLM provider, versioning, analytics)

commit f69da5d4110b802fff47901743fd7737ce9d74c3
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:23:39 2025 +0100

    DOCS: add decision support document with product areas, implications, complexity, and trade-offs

commit 74bd8a117d728ee495ee6ed8d4e61c9dd77a8163
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 11:16:10 2025 +0100

    DOCS: add Phase 2 PRD for Form System implementation

commit 89de36286e55a43f1b0130cd0c665746b10c2f3b
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:57:32 2025 +0100

    DOCS: add PRDs and integration documents

commit f34f9a905a1c58a9c66d237f06778ffb711c6e82
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:57:32 2025 +0100

    FEAT(ui): add TechChip component and examples under src/components

commit 87109da0bb998f8d7361ddc90f1ead7bcb25d319
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:51:55 2025 +0100

    FEAT(web): scaffold SvelteKit app with Tailwind and basic routes

commit 74c1271e6b38ccc970dc02e00249548a00eed755
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:51:34 2025 +0100

    DOCS: expand README with architecture, setup and testing sections

commit 7d3a6247af342093d4d95be5fd1dfa09230ce5a7
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:51:34 2025 +0100

    CHORE(devops): add docker-compose and CI workflow for migrations and linting

commit 82fba6d69b6a754892f8c7fb3effa9febb034807
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:51:34 2025 +0100

    FEAT(db): add Alembic setup and initial core tables migration

commit 6e79c8fad2ce8d1a0467bce2bfa0f80c28d5cc28
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:51:34 2025 +0100

    FEAT(api): scaffold FastAPI backend with health endpoint and CORS

commit a549c6056311ac5707932babcaa15c2c45650f09
Author: rickhallett <rickhallett@icloud.com>
Date:   Sat Sep 13 10:51:34 2025 +0100

    DOCS: add AGENTS.md contributor guidelines

---