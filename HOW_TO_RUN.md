# iSPARC — How to Run (Local Development)

* [ ] Step-by-step guide to run the full iSPARC stack (SvelteKit web + Go Fiber API + PostgreSQL) on a local machine, validate it, and commit/push the work on this branch. Written against branch `local-validation` (based on `main` @ `06ebd7c`, identical to `upstream/main`).

---

## 1. Prerequisites

| Tool           | Minimum version      | Notes                                                                             |
| -------------- | -------------------- | --------------------------------------------------------------------------------- |
| Node.js        | v20+ (tested on v24) | with npm 10+                                                                      |
| Go             | 1.24+                | `go.mod` targets 1.25; Go auto-downloads the right toolchain (`GOTOOLCHAIN=auto`) |
| Docker Desktop | any recent           | only used for PostgreSQL locally                                                  |
| Git            | any recent           |                                                                                   |

* [ ]

* `gofmt -l .` will list **every** file locally. This is a CRLF line-ending artifact of `core.autocrlf=true`, not real formatting drift — CI (Linux) is the source of truth. `gofmt -d <file>` showing only whitespace confirms it.
* `go test -race` fails locally with `cc1.exe: 64-bit mode not compiled in` unless you install a 64-bit MinGW gcc. Run tests **without** `-race` locally; CI runs the race detector.
* If port 3000/8080/5432 is already taken, check for leftover containers from other checkouts: `docker ps` → `docker compose stop` inside whichever project owns them (e.g. the older `ispark-v1` stack from `B:\Public Repository\SIPS\ispark_fork\ispark`).

---

## 2. One-time setup

### 2.1 Install dependencies

```powershell
# repo root — husky/lint-staged hooks
npm install

# frontend
cd web
npm install
cd ..

# backend (optional; go run resolves modules automatically)
cd api
go mod download
cd ..
```

### 2.2 Create env files (gitignored — never commit these)

`api/.env`:

```ini
PORT=8080

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=user
DB_NAME=isparc
DB_SSLMODE=disable

# Generate your own, e.g.:  node -e "console.log(require('crypto').randomBytes(48).toString('base64url'))"
JWT_SECRET=<paste-a-long-random-string>
JWT_REFRESH_SECRET=<paste-a-different-long-random-string>

# Leave SMTP empty locally — OTPs are printed to the API console instead
SMTP_HOST=
SMTP_PORT=
SMTP_USER=
SMTP_PASS=
SMTP_SENDER=

SEED_DEV_DATA=true
```

`web/.env`:

```ini
PUBLIC_API_BASE_URL=http://localhost:8080
```

---

## 3. Run the stack (three terminals)

### Terminal 1 — database

```powershell
docker compose up -d db
```

Verify: `docker ps` should show `ispark-db-1 … 0.0.0.0:5432->5432/tcp`. If the port mapping is missing (container shows just `5432/tcp`), the container was created while the port was busy — recreate it: `docker compose up -d --force-recreate db`.

### Terminal 2 — API

```powershell
cd api
go run .
```

On first boot it auto-migrates the schema and (because `SEED_DEV_DATA=true`) seeds demo data:
**8 students, 3 admins, 2 tracks, 7 activities, enrollments, certificates, settings, announcements.**
Wait for: `Starting server on port 8080...`

### Terminal 3 — web

```powershell
cd web
npm run dev
```

Open **[http://localhost:5173](http://localhost:5173)**.

### Demo credentials (every account's password is `Pass@123`)

| Portal                              | Login                   | Notes                              |
| ----------------------------------- | ----------------------- | ---------------------------------- |
| Student (`/login`)                  | `rahul.sharma@iips.edu` | IT2K24, most data                  |
| Student                             | `sneha.kumar@iips.edu`  | IT2K24                             |
| Student                             | `arjun.desai@iips.edu`  | IT2K24, has a rejected certificate |
| Student                             | `priya.nair@iips.edu`   | IT2K25                             |
| Student                             | `vikram.singh@iips.edu` | IT2K24, no activity                |
| Admin (`/admin-portal`)             | `admin`                 | batch IT2K24                       |
| Admin                               | `admin2`                | batch IT2K25                       |
| Super admin (`/super-admin-portal`) | `superadmin`            | whole platform                     |

Each login page also has a collapsible **Dev: demo login credentials** panel.

OTPs (registration / forgot-password) are **printed in the API console** (Terminal 2) when SMTP is unconfigured. You can also read them from the DB:

```powershell
docker exec ispark-db-1 psql -U postgres -d isparc -t -A -c "select code from otps where email='YOUR_EMAIL' order by created_at desc limit 1;"
```

### Alternative: full Docker stack

```powershell
docker compose up --build
```

Web on **[http://localhost:3000](http://localhost:3000)**, API on 8080, pgAdmin on 5050 (`admin@admin.com` / `user`). Uploads persist in the `api_uploads` volume.

---

## 4. Validate (mirror of CI)

Frontend (from `web/`):

```powershell
npm run format:check   # prettier
npm run lint:eslint    # eslint
npm run check          # svelte-check (TS diagnostics)
npm run build          # production build (adapter-node)
npm audit --audit-level=high
```

Backend (from `api/`):

```powershell
go vet ./...
go build ./...
go test -count=1 ./...       # unit tests run on in-memory SQLite, no DB needed
golangci-lint run ./...      # CI uses v2.12.1
```

All of the above pass on this branch as of 2026-07-22 (see `test.md` §0 for the recorded results).

### Automated API regression suite

With db + API running:

```powershell
node scripts/api-regression.mjs
```

63 end-to-end checks across auth (captcha/login/register/OTP/refresh/forgot-reset), the student module (dashboard, activities, enroll, certificates upload/download + security, leaderboard, marksheet, profile), RBAC boundaries, the admin module (batch scoping, cross-batch denial), and the super-admin platform APIs (users, settings, tracks CRUD, announcements CRUD/publish). Current result: **63/63 pass** — in local-disk storage mode **and** in Vercel Blob mode (run against a protocol-accurate local mock of the Blob API). The storage layer also has Go unit tests (`api/storage/storage_test.go`).

Re-runs are safe: it registers a fresh throwaway student each run. Point it at another deployment with `API_BASE_URL=https://... node scripts/api-regression.mjs` (the register/reset OTP checks need DB access and will fail without the local docker container — everything else works remotely).

---

## 5. Committing this branch's changes

Nothing has been committed for you — review, then run these manually.

```powershell
git status                       # review what changed
git diff                         # review content

# stage everything on this branch: docs, regression suite, storage layer,
# deployment env support, adapter switch, compose pgadmin opt-in, Render blueprint
git add HOW_TO_RUN.md test.md .gitignore docker-compose.yml render.yaml `
  scripts/api-regression.mjs `
  api/storage/ api/main.go api/config/database.go api/config/seed.go `
  api/controllers/student_dashboard_controller.go api/.env.example `
  web/vite.config.ts web/package.json web/package-lock.json

git commit -m "feat: deployment-ready test branch — Vercel Blob storage, env CORS, Supabase pooler support, Vercel adapter, run guide + test plan + API regression suite"

# publish the branch to your fork
git push -u origin local-validation
```

To later open a PR against upstream:

```powershell
gh pr create --repo iips-oss/ispark --base main --head ARPANPATRA111:local-validation `
  --title "Validation: run guide, manual test plan, API regression suite" `
  --body "Adds HOW_TO_RUN.md, test.md (manual regression plan), and scripts/api-regression.mjs (63 automated API checks)."
```

To keep the branch synced with upstream while it lives on the fork:

```powershell
git fetch upstream
git merge upstream/main          # or: git rebase upstream/main
```

---

## 6. Deployment (Supabase + Vercel + Vercel Blob)

The code on this branch is **deployment-ready** for the test environment. What changed to make that true:

| Concern                                    | Status                                                                                                                                                                                                                                                                                                    |
| ------------------------------------------ | --------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| Certificate files on ephemeral filesystems | **Solved** — `api/storage/` abstraction: local disk by default, **Vercel Blob** automatically when `BLOB_READ_WRITE_TOKEN` is set. Downloads stream through the API's auth check, so use a **private** Blob store. Verified end-to-end against a protocol-accurate mock (63/63 regression in blob mode). |
| Web hosting on Vercel                      | **Solved** — `web/vite.config.ts` picks `@sveltejs/adapter-vercel` when built on Vercel (`VERCEL` env), `adapter-node` everywhere else. Local/CI/Docker builds are unchanged.                                                                                                                            |
| CORS for the deployed web origin           | **Solved** — set `ALLOWED_ORIGINS=https://<your-app>.vercel.app` on the API (comma-separated, appended to the localhost defaults).                                                                                                                                                                       |
| Supabase transaction pooler (port 6543)    | **Solved** — set `DB_PREFER_SIMPLE_PROTOCOL=true` (disables pgx prepared statements). Not needed on the session pooler / direct connection (port 5432).                                                                                                                                                  |
| Seeding                                    | Already env-driven — first boot with `SEED_DEV_DATA=true`, then set `false`. Idempotent, and seed files upload to Blob when Blob is active.                                                                                                                                                              |

> Windows note: `VERCEL=1 npm run build` locally fails with a symlink `EPERM` unless Windows Developer Mode is on. Irrelevant for real deploys — Vercel builds on Linux.

**Chosen test stack: Vercel (web) + Render (API) + Supabase (DB) + Vercel Blob (files) — all free tiers.** Deploy in this order (each step feeds the next):

### 6.1 Supabase (PostgreSQL) — DONE ✔

Organisation **`ispark-deploy`**, project **`ispark-render`** (Mumbai). The project ref, host and password are in the gitignored `api/.env.supabase`; open the project from [supabase.com/dashboard](https://supabase.com/dashboard). If the dashboard looks empty, switch organisation with the picker at the top-left; it defaults to whichever org you used last. **Schema and demo data are already provisioned** and the 63-check regression suite passed against this database end to end.

> An earlier project `ispark-test` (org `ispark-oss`) is also still active. Nothing points at it any more — pause or delete it in the dashboard to free a free-tier slot (max 2 active projects per account).

All credentials live in **`api/.env.supabase`** (gitignored; template committed as `api/.env.supabase.example`). That one file holds the Supabase connection, the Blob token, the Brevo SMTP settings, and the CORS origin — i.e. exactly what Render needs.

Run the API against the cloud stack instead of the local Docker database:

```powershell
cd api
$env:ENV_FILE=".env.supabase"; go run .
# back to local:  Remove-Item Env:\ENV_FILE
```

Connection summary (session pooler, IPv4): host `aws-<n>-<region>.pooler.supabase.com`, port `5432`, user `postgres.<project-ref>`, database `postgres`, `sslmode=require` — copy the exact values from Supabase → **Connect** → *Session pooler*. Note the host's `aws-N` prefix differs per project; `aws-0` does not necessarily serve yours.

> The Supabase **CLI** and the Supabase **dashboard** authenticate separately. Signing into a different account in the browser does not move the CLI — run `npx supabase login` again to re-point it, and confirm with `npx supabase orgs list`.

> Free-tier caveat: Supabase **pauses projects after ~7 days of inactivity** — restore from the dashboard, or have someone log into the test portal weekly.

### 6.0 Live URLs

| Piece        | URL                                                    |
| ------------ | ------------------------------------------------------ |
| Web (Vercel) | <https://ispark-roan.vercel.app>                       |
| API (Render) | <https://ispark-api.onrender.com> — health: `/health` |
| Database     | Supabase `ispark-render` (Mumbai)                      |
| Files        | Vercel Blob, private store `ispark-certificates`       |

Verified against the live stack: `/health` OK, CORS allows the Vercel origin, the frontend resolves to the Render API, and `scripts/api-regression.mjs` passes **63/63**.

```powershell
# re-run the full suite against the deployment at any time
$env:API_BASE_URL="https://ispark-api.onrender.com"
$env:OTP_DB_URL=(Select-String -Path api\.env.supabase -Pattern '^DATABASE_URL=(.+)$').Matches[0].Groups[1].Value
$env:TEST_EMAIL_BASE="you@example.com"   # real inbox: avoids bounces, since SMTP is live
node scripts/api-regression.mjs
```

### 6.0.2 Keeping the free tiers awake

Two sleep policies apply: Render spins the API down after **15 minutes** without an inbound HTTP request, and Supabase pauses a project after **7 days** without database activity. Both are covered by two independent mechanisms, so neither is a single point of failure.

**The key idea:** `GET /health/db` runs a real `SELECT` against the database. One request to it therefore resets Render's idle timer *and* registers as Supabase activity — a single ping solves both problems. (`/health` stays database-free so a database outage can never fail Render's own health check and take the API down.)

| Mechanism                                                 | Runs where                  | Covers            |
| --------------------------------------------------------- | --------------------------- | ----------------- |
| Uptime monitor calling `/health/db` every 10 min          | external (UptimeRobot etc.) | Render + Supabase |
| `keepalive-ping-api` pg_cron job, every 10 min via pg_net | inside Supabase             | Render + Supabase |
| `keepalive-write` / `keepalive-delete` pg_cron jobs       | inside Supabase             | Supabase          |

Apply the database side with:

```powershell
node scripts/db.mjs --file scripts/sql/keepalive.sql
```

That creates `ops.keepalive` (in a non-public schema, so PostgREST cannot expose it) and schedules three jobs: a heartbeat row written at 10:00 on days 1/6/11/16/21/26, deleted at 08:00 the next day, and an HTTP call to `/health/db` every 10 minutes. Schedules are **UTC** — subtract 5h30m from an IST time.

Inspect or change them:

```powershell
node scripts/db.mjs "select jobname, schedule, active from cron.job"
node scripts/db.mjs "select jobname, status, start_time from cron.job_run_details order by start_time desc limit 10"
node scripts/db.mjs "select status_code, created from net._http_response order by id desc limit 5"
```

> **Render does not detect or block keep-alive traffic.** Its rule is purely "no inbound request for 15 minutes". Staggering monitors to look human is unnecessary; a second monitor is only worth adding for redundancy if the first lapses.

> **Watch the free-tier budget.** Keeping the API awake 24/7 costs ~730 of Render's 750 free instance-hours per month, leaving almost no headroom for a second service. If you only need it responsive during working hours, pause the monitor overnight (or narrow the pg_cron schedule to, say, `*/10 3-18 * * *` UTC = 08:30–23:30 IST) and roughly halve the usage.

### 6.0.1 Database access and security

Query either database without Docker or psql:

```powershell
node scripts/db.mjs "select count(*) from students"          # cloud (api/.env.supabase)
node scripts/db.mjs --local "select count(*) from students"  # local  (api/.env)
node scripts/db.mjs --file scripts/sql/harden-public-schema.sql
```

**Supabase exposes every table in the `public` schema over PostgREST**, authorised by the `anon` key that is public by design. This project does not use Supabase Auth or PostgREST — the Go API owns all access control — so those roles are pure attack surface. `scripts/sql/harden-public-schema.sql` enables RLS (no policies) and revokes `anon`/`authenticated` privileges; the API is unaffected because the table owner bypasses RLS. Verified: the anon key now receives `42501 insufficient_privilege`.

> **Re-run that script after any schema change.** `AutoMigrate` creates new tables without RLS, which re-opens the hole. The script is idempotent.

> **Supabase → Authentication → Users will always be empty.** Students and admins live in our own `students`/`admins` tables with our own JWTs; Supabase Auth is not used. Look at **Table Editor → students** instead.

### 6.2 Render (API) — no CLI needed, and that's fine

The Render CLI is **not required at all** (it's mostly for log tailing/SSH anyway — nothing in this workflow needs it). Render's native deploy path is Git + Blueprint, driven entirely from the dashboard:

1. Push this branch to your fork (`git push -u origin local-validation`)
2. [dashboard.render.com](https://dashboard.render.com) → **New + → Blueprint** → connect GitHub → pick `ARPANPATRA111/ispark` → branch `local-validation`
3. Render reads **`render.yaml`** (repo root, already written): Docker service `ispark-api` from `api/Dockerfile`, free plan, Singapore region, health check on `/health`, JWT secrets auto-generated
4. The dashboard prompts for the `sync: false` values: Supabase `DB_HOST/DB_USER/DB_PASSWORD`, `BLOB_READ_WRITE_TOKEN` (copy from gitignored `web/.env.local`), SMTP (optional — leave blank and OTPs land in Render's log stream). `ALLOWED_ORIGINS` is pre-filled with the live Vercel URL
5. Deploy → note the URL, e.g. `https://ispark-api.onrender.com` → check `https://ispark-api.onrender.com/health`
6. After first successful boot with data, set `SEED_DEV_DATA=false` in the service's Environment tab

> Free-tier caveat: the service **sleeps after ~15 min idle**; the next request cold-starts in ~1 min. Free allowance (750 h/month) covers one always-on service, so an [UptimeRobot](https://uptimerobot.com) (free) monitor pinging `/health` every 10 min keeps it awake during the testing window.

### 6.3 Vercel (web) — DEPLOYED ✔

Live at **<https://ispark-roan.vercel.app>** (project `ispark`, scope `arpanpatra111s-projects`, CLI deploy from local files — no GitHub connection, so you can still wire up the git-based deployment from branch `local-validation` in the dashboard whenever you want; the project and env vars are already there).

Project env vars already set: `PUBLIC_API_BASE_URL` (`https://ispark-api.onrender.com` for production/preview, localhost for development) and `BLOB_READ_WRITE_TOKEN` (all environments). If Render assigns a different URL than `ispark-api.onrender.com`, update the env and redeploy:

```powershell
cd web
vercel env rm PUBLIC_API_BASE_URL production   # then re-add with the real URL
"https://<real-url>" | vercel env add PUBLIC_API_BASE_URL production
vercel deploy --prod --yes --archive=tgz
```

### 6.4 Vercel Blob (certificate storage) — DONE ✔

Private store **`ispark-certificates`** is created and **connected to the `ispark` Vercel project**. The long-lived `BLOB_READ_WRITE_TOKEN` lives in the gitignored **`web/.env.local`** (and in the Vercel project's env). One task left:

* [ ] Copy the token value from `web/.env.local` into **Render's** `BLOB_READ_WRITE_TOKEN` env var when the Blueprint prompts for it — the Go API is the only thing that talks to Blob; the browser never does. Never commit it.

### 6.α Free hosting options compared (for this test deployment)

| Option (API host)     | Free tier                                                                     | Verdict                                                                           |
| --------------------- | ----------------------------------------------------------------------------- | --------------------------------------------------------------------------------- |
| **Render** (chosen)   | 750 h/mo, 512 MB, sleeps after 15 min idle, no card needed, Blueprint deploys | **Best free fit** — Docker support, health checks, secret generation             |
| Railway               | $5 one-time trial credit, then paid                                           | Original plan target; fine later for production, not free for ongoing testing     |
| Fly.io                | No real free tier anymore (pay-as-you-go, card required)                      | Good product, not "free"                                                          |
| Koyeb                 | 1 free web service, scale-to-zero                                             | Workable backup if Render misbehaves                                              |
| Google Cloud Run      | Generous free tier but requires a credit card + GCP setup                     | Overkill for a test deployment                                                    |
| Vercel functions (Go) | Free                                                                          | Would require rewriting the Fiber server into serverless handlers — not worth it |

Web on **Vercel Hobby** (free) and DB on **Supabase Free** (500 MB, 2 projects) are the clear choices; no real competition at free tier. For production later: Railway or a paid Render/Fly instance for the API, same Vercel/Supabase paid tiers as usage demands.

### 6.5 Remaining production to-dos

* [ ] Real SMTP credentials (Brevo) — without them OTPs only go to server logs
* [ ] Strong `JWT_SECRET`/`JWT_REFRESH_SECRET` (captcha falls back to `default_secret` if unset!)
* [ ] Remove/feature-flag `web/src/lib/DevCredentials.svelte` (imported on all three login pages) before real users see it
* [ ] `SEED_DEV_DATA=false` after first boot
* [ ] Certificate approve/reject backend still missing upstream (PR #102) — the credit loop can't complete on any deployment until it lands

---

## 7. Layout of this repo

```text
root/
├── web/          SvelteKit 2 + Svelte 5 + TS + Tailwind v4
│                 (adapter-node locally/CI, adapter-vercel on Vercel)
├── api/          Go Fiber v2 + GORM; unit tests use in-memory SQLite
│   └── storage/  file storage abstraction: local disk ⇄ Vercel Blob
├── scripts/      api-regression.mjs — automated API regression suite (63 checks)
├── mks/          (gitignored) local project context for maintainers/LLMs
├── docker-compose.yml   web:3000 api:8080 db:5432 (pgadmin opt-in: --profile tools)
├── render.yaml          Render Blueprint — dashboard-based API deploy, no CLI needed
├── HOW_TO_RUN.md        this file
└── test.md              manual test plan (assign sections to contributors)
```

Docker hygiene: only the `db` container is required for development (`docker compose up -d db`). pgAdmin is behind the `tools` profile so it is never pulled by default. When you are done working: `docker compose stop`. To reclaim space periodically: `docker builder prune` and `docker volume prune` (never `--all`, and never prune named volumes — they hold database data).
