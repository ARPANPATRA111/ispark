# iSPARC — Manual Test Plan (Regression / Pre-Production)

Purpose: human validation of everything the automated checks cannot see (visual state, UX flows, browser behavior, real email, deployed infra). Each numbered section below is independently assignable to one contributor. Fill in the **Tester** line, tick the checkboxes, and record any failure with steps + screenshot in the *Findings* line of that section.

**Test environment — LIVE.** Web: **<https://ispark-roan.vercel.app>** · API: **<https://ispark-api.onrender.com>** (health check: `/health`). Or run locally per `HOW_TO_RUN.md`.

> **First request of the day is slow.** The API sleeps after ~15 minutes of inactivity on Render's free plan, so the first login can take up to a minute while it wakes. That is expected — not a bug. Please don't file it.

**Credentials:** all seeded accounts use password `Pass@123`.
Students: `rahul.sharma@iips.edu` (IT2K24, most data), `sneha.kumar@iips.edu`, `arjun.desai@iips.edu` (has a rejection), `priya.nair@iips.edu` (IT2K25), `vikram.singh@iips.edu` (no activity).
Admins: `admin` (IT2K24), `admin2` (IT2K25). Super admin: `superadmin`.

**OTPs:** locally they print in the API console. On staging they arrive by email only if SMTP is configured — otherwise ask the maintainer to read them from the DB.

---

## 0. Automated baseline (already executed — for reference, not re-testing)

Recorded 2026-07-22 on branch `local-validation` (`main` @ `06ebd7c`):

| Check                                                                                             | Result                                                                                               |
| ------------------------------------------------------------------------------------------------- | ---------------------------------------------------------------------------------------------------- |
| Prettier / ESLint / svelte-check                                                                  | pass / pass / 0 errors                                                                               |
| Frontend production build (adapter-node)                                                          | pass                                                                                                 |
| `npm audit --audit-level=high`                                                                    | pass (3 low-sev via SvelteKit `cookie` chain — known, tracked)                                      |
| `go vet`, `go build`, `golangci-lint` (v2.12.1)                                                   | pass, 0 issues                                                                                       |
| `go test ./...` (controllers + storage suites)                                                    | pass                                                                                                 |
| API regression suite `scripts/api-regression.mjs` — local-disk storage                           | **63/63 pass**                                                                                       |
| API regression suite — **Vercel Blob mode** (against protocol mock, incl. private download auth) | **63/63 pass**                                                                                       |
| API regression suite — **against the live deployment** (Render API + Supabase + Vercel Blob)     | **63/63 pass**                                                                                       |
| All 9 web routes SSR 200                                                                          | pass                                                                                                 |
| Vercel adapter auto-selection (`VERCEL=1` build picks adapter-vercel)                             | pass (full local build needs Windows Developer Mode for symlinks; real builds run on Vercel's Linux) |

The automated suite already covers API happy paths, auth failures, RBAC boundaries, batch scoping, file-type sniffing, and cross-student certificate access. **Human testing should focus on the UI wiring, visuals, and the flows below.**

### ⚠ Views that are KNOWN mock/prototype (no backend exists yet — do not file bugs for fake data)

| View                                                        | Upstream status                                                                                         |
| ----------------------------------------------------------- | ------------------------------------------------------------------------------------------------------- |
| Admin → Dashboard stat cards                               | not wired (PR #88/#89 attempt open/closed)                                                              |
| Admin → Certificate Verification                           | **no backend endpoints at all** (PR #102 open) — the product's approve/reject loop cannot complete yet |
| Admin → Activity Monitoring                                | mock                                                                                                    |
| Admin → Batch Analytics                                    | mock                                                                                                    |
| Super admin → Activity Management                          | local-state only, changes vanish on refresh (API PR #91 open)                                           |
| Super admin → Reports Center                               | mock (API PR #98 open)                                                                                  |
| Super admin dashboard sub-labels ("+32 this semester" etc.) | fabricated deltas on real numbers                                                                       |

For these, the test is only: *page renders, doesn't crash, and is visibly a prototype*. Everything else in this plan is expected to be fully functional.

---

## A. Public pages & authentication

### A1. Landing page & navigation

**Tester:** __________

* [ ] `/` renders hero, tracks, timeline, outcomes, about, footer; no broken images or console errors
* [ ] Navbar links scroll/route correctly; active section stays highlighted while scrolling
* [ ] "Login" / "Register" buttons route to `/login` and `/register`
* [ ] Page is usable at 375px, 768px, 1440px widths

**Findings:** __________

### A2. Student registration + OTP verification

**Tester:** __________

* [ ] Register a brand-new student (unique email + roll no + enrollment no) → success message, OTP issued
* [ ] Submit wrong OTP → clear error, account stays unverified
* [ ] Submit correct OTP → lands in portal logged in
* [ ] Try registering again with the same email → rejected with readable error
* [ ] Try mismatched password/confirm-password → rejected client- or server-side
* [ ] Weak password (e.g. `abc`) → rejected with the password policy message
* [ ] Login attempt on an unverified account → blocked and a fresh OTP is sent

**Findings:** __________

### A3. Student login, captcha, forgot/reset password

**Tester:** __________

* [ ] Login with `rahul.sharma@iips.edu` + correct captcha → portal loads
* [ ] Wrong password → "Invalid credentials", no session
* [ ] Wrong captcha answer → error; a new captcha can be requested
* [ ] Refresh: after login, reloading `/portal` keeps you logged in (refresh token cookie)
* [ ] Logout → back to login; `/portal` no longer shows data
* [ ] Forgot password: request OTP for a seeded student → OTP arrives (console/email) → reset succeeds → old password stops working, new one works
  * ⚠ if you reset a shared demo account, set the password back to `Pass@123` afterwards
* [ ] Reset with expired/wrong OTP → clear error

**Findings:** __________

## B. Student portal (`/portal`)

### B1. Dashboard & stats

**Tester:** __________

* [ ] Stat cards (credits, enrollments, certificates, rank) show **real numbers** consistent with the account used (e.g. `vikram.singh` ≈ zeros, `rahul.sharma` non-zero)
* [ ] Recent activity list matches the account's actual enrollments/uploads
* [ ] No hardcoded dates/academic year on the dashboard (compare with current date)
* [ ] Numbers update after you enroll or upload (revisit after B3/B4)

**Findings:** __________

### B2. My profile & change password

**Tester:** __________

* [ ] Profile shows the seeded data for the logged-in student
* [ ] Edit contact/details → save → refresh → change persisted
* [ ] Change password with wrong current password → rejected
* [ ] Change password properly → old rejected, new works (restore `Pass@123` afterwards)

**Findings:** __________

### B3. Browse activities & enrollment

**Tester:** __________

* [ ] Catalogue lists 7 seeded activities with category grouping (no duplicate categories from casing, e.g. `TECHNICAL` vs `Technical`)
* [ ] Enroll in an open activity → success feedback; it appears under My Enrollments
* [ ] Enrolling twice in the same activity → blocked with readable error
* [ ] My Enrollments statuses (enrolled/completed) look consistent with seeded data

**Findings:** __________

### B4. Upload certificate & My certificates

**Tester:** __________

* [ ] Upload a real PDF with all fields → appears in My Certificates as **Pending**
* [ ] Upload a non-PDF renamed to `.pdf` → rejected (server sniffs content)
* [ ] Upload > 5 MB file → rejected with size message
* [ ] Download your own uploaded certificate → correct file opens
* [ ] `arjun.desai` account: rejected certificate is visibly marked rejected
* [ ] Status filters/tabs on My Certificates work

**Findings:** __________

### B5. Credits & progress + Extracurricular marksheet

**Tester:** __________

* [ ] Credits/progress figures match the same account's dashboard numbers (no contradiction between pages)
* [ ] Academic-year label shown is the **current** academic year, not hardcoded
* [ ] Marksheet rows correspond to real (seeded/approved) activities of the account
* [ ] Note: the graduation target (e.g. `/200`) is a placeholder policy — record what is displayed

**Findings:** __________

### B6. Leaderboard

**Tester:** __________

* [ ] Leaderboard shows ranked seeded students with points; logged-in student's row highlighted
* [ ] Category champions row/cards populated
* [ ] Year filter changes results sensibly
* [ ] Dashboard credits vs leaderboard points: check the same student for contradictions (all-time vs academic-year windows — historically buggy, see mks audit §4.1)

**Findings:** __________

## C. Admin portal (`/admin-portal`)

### C1. Admin login & forced password change

**Tester:** __________

* [ ] `admin` / `Pass@123` logs in
* [ ] A newly created admin (from super admin User Management) is forced through the change-password screen on first login and can't skip it
* [ ] Wrong credentials → error; student credentials on admin portal → rejected

**Findings:** __________

### C2. Student management & student detail (wired)

**Tester:** __________

* [ ] `admin` sees exactly the IT2K24 students (5 seeded); `admin2` sees IT2K25 (3)
* [ ] Search/filter in the student list works
* [ ] Opening a student shows real certificates + enrollments for that student
* [ ] URL tampering: as `admin2`, try an IT2K24 student's detail URL directly → denied

**Findings:** __________

### C3. Admin profile (view + edit)

**Tester:** __________

* [ ] Profile shows the logged-in admin's real data (name, batch)
* [ ] Edit + save persists after refresh
* [ ] Change password flow works (restore `Pass@123` afterwards)

**Findings:** __________

### C4. Prototype views render safely (Dashboard stats, Certificate Verification, Activity Monitoring, Batch Analytics)

**Tester:** __________

* [ ] Each view opens without crashing or console errors
* [ ] Confirm they still show **mock** data (see table in §0) — record anything that has silently become half-wired
* [ ] Certificate Verification: confirm there is **no way** to actually approve/reject (backend missing, PR #102)

**Findings:** __________

## D. Super admin portal (`/super-admin-portal`)

### D1. Login & dashboard

**Tester:** __________

* [ ] `superadmin` / `Pass@123` logs in; a plain `admin` account is **rejected** from the super admin portal
* [ ] The four stat cards show real platform numbers (compare with seeded totals: 8+ students, 3 admins, 7 activities)
* [ ] Note the sub-labels ("+32 this semester" …) are fabricated — record, don't fail

**Findings:** __________

### D2. User management (wired)

**Tester:** __________

* [ ] User registry lists seeded students + admins with roles
* [ ] Create a new Admin user → appears in list; temp password shown once
* [ ] Log in as that new admin (forced password change — feeds C1)
* [ ] Delete the test user → gone after refresh
* [ ] Try deleting your own `superadmin` account → blocked

**Findings:** __________

### D3. Track management (wired)

**Tester:** __________

* [ ] Track list + stats show seeded tracks with real activity counts
* [ ] Create a track → persists after refresh; duplicate name → conflict error
* [ ] Edit description/status (Active/Inactive) → persists
* [ ] Delete a test track → gone after refresh; check behavior when a track still has activities (record what happens)

**Findings:** __________

### D4. Announcement management (wired)

**Tester:** __________

* [ ] Seeded announcements list with correct status chips (active/draft/expired/scheduled)
* [ ] Create draft → publish → status changes and persists after refresh
* [ ] Edit + delete announcement persist after refresh
* [ ] Date validation: expiry before publish date → rejected

**Findings:** __________

### D5. System settings (wired)

**Tester:** __________

* [ ] Settings load grouped by category (Academic Year etc.)
* [ ] Change a setting value → save → refresh → persisted
* [ ] Invalid value (if validation exists) → readable error; record which settings actually affect app behavior (many are display-only for now)

**Findings:** __________

### D6. Prototype views render safely (Activity Management, Reports Center)

**Tester:** __________

* [ ] Activity Management: CRUD *appears* to work but changes vanish on refresh — confirm, record
* [ ] Reports Center: renders, all data mock, export buttons don't error the page

**Findings:** __________

## E. Cross-cutting

### E1. Security & RBAC (manual spot-checks on top of automated coverage)

**Tester:** __________

* [ ] Logged out: hitting `/portal`, `/admin-portal/dashboard`, `/super-admin-portal/dashboard` directly by URL → no data renders / redirected (note: client-side guard only; API returns 401 regardless)
* [ ] Student token in devtools cannot fetch `/api/admin/*` (403)
* [ ] Session expiry: leave a tab idle past token expiry → app recovers via refresh or asks to re-login rather than half-broken UI
* [ ] No secrets/OTPs/passwords visible in browser devtools network responses beyond what the page needs

**Findings:** __________

### E2. Responsiveness & cross-browser

**Tester:** __________

* [ ] All portals usable at 375px (phone), 768px (tablet), 1440px: nav collapses, tables scroll, modals fit
* [ ] Chrome + Firefox + one mobile browser: no layout breakage or dead buttons
* [ ] Dark-mode/OS theme doesn't render text unreadable (if applicable)

**Findings:** __________

### E3. Staging deployment checks (fill in once Supabase + Vercel are live)

**Tester:** __________

Already verified by the maintainer: the API answers `/health`, CORS allows the Vercel origin, the frontend is wired to the API, and the 63-check suite passes against the live stack. The items below need a human.

* [ ] Registration OTP arrives by **real email** to a genuine inbox (Brevo is live on the deployment) — check spam too
* [ ] Cold start: after ~15 min idle, the first request takes up to a minute, then the app is responsive
* [ ] Certificate upload lands in the **Vercel Blob store** (check the store dashboard, `certificates/` prefix) and downloads back correctly via the portal
* [ ] **Redeploy the API** → previously uploaded certificate still downloads (Blob storage survives; local-disk mode would not)
* [ ] Blob store is **Private**: opening a raw blob URL directly in an incognito tab must NOT serve the file
* [ ] Seed data present exactly once (no duplicates from multiple boots); seeded certificate PDFs download correctly (they live in Blob too)
* [ ] HTTPS everywhere; refresh-token cookie works cross-site (SameSite/secure attributes)

**Findings:** __________

---

## Sign-off

| Section | Tester | Date | Pass/Fail |
| ------- | ------ | ---- | --------- |
| A1      |        |      |           |
| A2      |        |      |           |
| A3      |        |      |           |
| B1      |        |      |           |
| B2      |        |      |           |
| B3      |        |      |           |
| B4      |        |      |           |
| B5      |        |      |           |
| B6      |        |      |           |
| C1      |        |      |           |
| C2      |        |      |           |
| C3      |        |      |           |
| C4      |        |      |           |
| D1      |        |      |           |
| D2      |        |      |           |
| D3      |        |      |           |
| D4      |        |      |           |
| D5      |        |      |           |
| D6      |        |      |           |
| E1      |        |      |           |
| E2      |        |      |           |
| E3      |        |      |           |
