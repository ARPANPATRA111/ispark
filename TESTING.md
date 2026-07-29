# iSPARC — Manual Test Suite

This is the full manual test suite for the iSPARC testing round. It is written for **human testers**. Every case is meant to be executed by a person clicking through the live site.

**Do not use an AI agent, script, or browser automation to produce these results.** Automated checks already cover the API surface (63 endpoint checks run on every deploy). The entire value of this round is what an automated check cannot see: confusing wording, a button that looks disabled but isn't, a layout that breaks on a phone, a number that is technically correct but misleading, a flow that works but feels wrong. Report what you actually observed, including "this was confusing" — that is a valid finding.

---

## 1. Environment

|          |                                            |
| -------- | ------------------------------------------ |
| Web      | <https://ispark-iips.vercel.app>           |
| API      | <https://ispark-api.onrender.com>          |
| Database | Supabase (shared — see the warning below) |

> **The first request of the day is slow.** The API sleeps after ~15 minutes of inactivity and takes up to a minute to wake. A slow *first* login is expected and is **not a bug**. A slow *second* request is a bug — report it.

> **This is a shared database.** Everyone tests against the same data, so anything you create or change is visible to other testers. Prefer creating your own records over editing someone else's, name them recognisably (e.g. `S4-test-activity`), and restore any password you change back to `Pass@123`.

### Browsers to cover

Between the whole team, every suite should be run at least once on each of: **Chrome (desktop)**, **Firefox (desktop)**, **Safari or Chrome on a real phone**. Note which you used at the top of your report.

---

## 2. Test accounts

Every account below uses the password **`Pass@123`**.

| Role        | Login                   | Portal                | Notes                                      |
| ----------- | ----------------------- | --------------------- | ------------------------------------------ |
| Student     | `rahul.sharma@iips.edu` | `/login`              | IT2K24                                     |
| Student     | `sneha.kumar@iips.edu`  | `/login`              | IT2K24                                     |
| Student     | `arjun.desai@iips.edu`  | `/login`              | IT2K24                                     |
| Student     | `vikram.singh@iips.edu` | `/login`              | IT2K24 — keep this one empty as a control |
| Student     | `priya.nair@iips.edu`   | `/login`              | IT2K25 — different batch                  |
| Admin       | `admin`                 | `/admin-portal`       | Scoped to batch **IT2K24**                 |
| Admin       | `admin2`                | `/admin-portal`       | Scoped to batch **IT2K25**                 |
| Super admin | `superadmin`            | `/super-admin-portal` | Whole platform                             |

Each login page also has a collapsible **Dev: demo credentials** panel.

### Starting state: an empty canvas

The database was deliberately cleared before this round. **Only the accounts above exist.** There are no activities, tracks, certificates, enrolments or announcements — you create everything yourself, which is the point: it exercises the real product flow from zero and means no seeded data can mask a bug.

| Entity                                                               | Count at start                     |
| -------------------------------------------------------------------- | ---------------------------------- |
| Students                                                             | 8 (5 in IT2K24, 3 in IT2K25)       |
| Admins                                                               | 3 (2 batch admins + 1 super admin) |
| Platform settings                                                    | 26                                 |
| Activities, tracks, certificates, enrolments, announcements, reports | **0**                              |

Two consequences to expect, which are **not bugs**:

* Every list starts empty. An empty state should be a friendly message, never a blank panel, a permanent spinner, or `undefined`/`NaN` — if you see those, that *is* a bug worth reporting.
* Dashboards, leaderboards and analytics read zero until data exists.

### Suite order matters

Because nothing exists yet, some suites need data another suite creates. Run them roughly in this order, or create what you need yourself:

```
S10 (super admin: create a track, then an activity)
   └─> S3  (student: browse + enrol in that activity)
          └─> S4  (student: upload a certificate)
                 └─> S8  (admin: verify/approve it)
                        └─> S5  (student: credits + marksheet PDF now have content)
                               └─> S6/S9/S11 (leaderboard, analytics, reports have data)
```

S1, S2, S7 and S12 need no pre-existing data and can start immediately.

> If you are blocked waiting on data, say so in your report rather than skipping the case — "could not test, no approved certificate existed yet" is useful information.

---

## 3. How to report a finding

Log one row per finding in your suite's results table, and raise anything Critical/High to the maintainer immediately rather than waiting until you finish.

**Include every time:** what you did (numbered steps), what you expected, what actually happened, the account and browser used, and a screenshot. If it involves an error, open DevTools (F12) → Console and Network, and include the failing request's status code.

### Severity

| Severity     | Meaning                                                       | Examples                                                                                    |
| ------------ | ------------------------------------------------------------- | ------------------------------------------------------------------------------------------- |
| **Critical** | Data loss, a security hole, or a core flow completely blocked | You can see another student's certificate; approvals silently do not save; login impossible |
| **High**     | A main feature is broken, no workaround                       | Upload always fails; marksheet exports blank                                                |
| **Medium**   | Feature works but is wrong or awkward                         | A count is off by one; validation message is misleading                                     |
| **Low**      | Cosmetic                                                      | Misalignment, typo, inconsistent capitalisation                                             |

### Please **do not** report these — they are known and accepted

* Slow first request after idle (free-tier cold start).
* The super-admin dashboard's small sub-labels ("+32 this semester") — known fabricated deltas, already on the fix list.
* The yellow "Dev: demo credentials" panel on login pages — intentional for this round.
* OTP emails landing in spam.

---

## 4. Suite assignments

Twelve suites. One contributor per suite is ideal; S12 should be done by someone who has already finished another suite.

| Suite | Area                                                | Assigned to |
| ----- | --------------------------------------------------- | ----------- |
| S1    | Public site & registration                          |             |
| S2    | Login, session & password recovery                  |             |
| S3    | Student — dashboard, activities, enrolments        |             |
| S4    | Student — certificate upload & management          |             |
| S5    | Student — credits, marksheet & PDF export          |             |
| S6    | Student — leaderboard & profile                    |             |
| S7    | Admin — access, profile & student management       |             |
| S8    | Admin — certificate verification (the credit loop) |             |
| S9    | Admin — activity monitoring & batch analytics      |             |
| S10   | Super admin — users, activities & tracks           |             |
| S11   | Super admin — announcements, settings & reports    |             |
| S12   | Cross-cutting — security, responsive, resilience   |             |

---

## S1 — Public site & registration

**Tester:** ______  **Browser:** ______  **Date:** ______

### S1.1 Landing page

1. Open the site logged out. Every section renders: hero, tracks, timeline, outcomes, about, footer.
2. No broken images, no overlapping text, no horizontal scrollbar.
3. Navbar links move to the right sections; the active link highlights while scrolling.
4. "Login" and "Register" reach `/login` and `/register`.
5. DevTools Console shows no red errors.

### S1.2 Registration — the happy path

1. Register with details you control (unique email, roll number and enrolment number).
2. Expect a clear success state and a prompt for the OTP.
3. **The email contains a 6-digit code, not a link.** Enter it to finish.
4. You land in the student portal, logged in, with an empty dashboard (0 credits, no certificates).

### S1.3 Registration — validation edge cases

Each of these must be **rejected with a message that says what to fix**:

| #   | Input                                                      | Expected                                                                                          |
| --- | ---------------------------------------------------------- | ------------------------------------------------------------------------------------------------- |
| a   | Email already registered (`rahul.sharma@iips.edu`)         | Rejected, tells you the email is taken                                                            |
| b   | Roll number already used (`IT2K24011`)                     | Rejected                                                                                          |
| c   | Enrolment number already used                              | Rejected                                                                                          |
| d   | Password and confirm-password differ                       | Rejected before submit                                                                            |
| e   | Weak password (`abc`)                                      | Rejected, states the policy                                                                       |
| f   | Malformed email (`test@`, `test.com`, `a b@c.com`)         | Rejected                                                                                          |
| g   | Every field blank, press submit                            | Cannot submit; required fields marked                                                             |
| h   | Semester out of range (0, 15, negative, letters)           | Rejected                                                                                          |
| i   | Contact number with letters, or 5 digits, or 20 digits     | Rejected                                                                                          |
| j   | Leading/trailing spaces in email (` a@b.com `)             | Either trimmed and accepted, or clearly rejected — **not** silently creating a duplicate account |
| k   | Very long name (200+ characters)                           | Either limited or handled without breaking layout                                                 |
| l   | Name with an accent or non-Latin script (`José`, `प्रिया`) | Accepted and displayed correctly afterwards                                                       |

### S1.4 OTP edge cases

1. Enter a **wrong** OTP → clear error, account stays unverified, you can retry.
2. Enter a **malformed** OTP (`12`, `abcdef`, empty) → rejected sensibly.
3. Request a **second** OTP, then try the **first** one → the older code should be refused (report which behaviour you see).
4. Submit the correct OTP **twice** (press back, resubmit) → no crash, no duplicate account.

### S1.5 Interruption

1. Register, but close the tab **before** entering the OTP.
2. Now go to `/login` and sign in with those credentials.
3. Expect: you are told the account is unverified, a fresh code is sent, and **an OTP box appears on the login page** for you to enter it. There should never be a message telling you a code was sent with nowhere to type it.

**Findings:**

| #   | Severity | Summary | Steps / expected / actual | Screenshot |
| --- | -------- | ------- | ------------------------- | ---------- |
|     |          |         |                           |            |

---

## S2 — Login, session & password recovery

**Tester:** ______  **Browser:** ______  **Date:** ______

### S2.1 Login

1. Sign in as `rahul.sharma@iips.edu` with the correct captcha answer → portal loads.
2. Wrong password → "Invalid credentials", no session created.
3. Unknown email → same generic message (it must **not** reveal whether the account exists).
4. Wrong captcha answer → rejected, and a **new** captcha appears.
5. Empty captcha → rejected.
6. Press the captcha **refresh** button → the question changes and remains solvable. Do this 5 times.
7. Leave the login page open for ~6 minutes, then submit → an expired-captcha error, and you can recover by refreshing it.

### S2.2 Session behaviour

1. After login, reload `/portal` → you stay logged in.
2. Open `/portal` in a second tab → still logged in.
3. Log out → returned to login; pressing **browser Back** must **not** show the portal with data.
4. After logging out, paste `/portal` directly → no student data is visible.
5. Log in, then clear `localStorage` in DevTools and reload → treated as logged out, no crash.

### S2.3 Forgot / reset password

Use `vikram.singh@iips.edu` (the empty account) so a password change affects nobody else. **Restore it to `Pass@123` when finished.**

1. Request a reset code → the email arrives.
2. Wrong code → clear error.
3. Correct code + a new valid password → success.
4. The **old** password no longer works; the **new** one does.
5. New password failing the policy → rejected with the rule stated.
6. New password and confirmation differing → rejected.
7. Request a reset for an email that does not exist → the response must **not** confirm whether the account exists.
8. Reuse an already-used reset code → refused.

### S2.4 Cross-portal access

1. Student credentials on `/admin-portal` → refused.
2. Admin credentials (`admin`) on the student `/login` → refused.
3. Admin credentials on `/super-admin-portal` → refused (this is a different privilege level).

**Findings:** *(same table format as S1)*

---

## S3 — Student: dashboard, activities & enrolments

**Tester:** ______  **Browser:** ______  **Date:** ______

Primary account: `rahul.sharma@iips.edu`. Keep `vikram.singh@iips.edu` untouched as an always-empty control for comparing empty states.

### S3.1 Dashboard accuracy

1. Stat cards (credits, enrolments, certificates, rank) show **real** numbers.
2. Cross-check: the certificate count must equal the rows under **My Certificates**.
3. Log in as `vikram.singh` (the untouched control) → the cards read zero with a sensible empty state, **not** blank boxes, "NaN", "undefined" or a spinner that never stops.
4. The academic year label shows the **current** year, not a hardcoded past one.
5. "Recent activity" matches things that account actually did.

### S3.2 Browse activities

1. Every activity created in S10 appears with title, category, credits and dates. (If none exist yet, confirm the empty state is friendly, then come back after S10.)
2. Category grouping shows no duplicate category caused by capitalisation (e.g. `TECHNICAL` and `Technical` as separate groups).
3. Search and filters return the right subset; clearing them restores the full list.
4. Search for something that cannot match (`zzzzz`) → a proper "no results" message.

### S3.3 Enrolment

1. Enrol in an activity → immediate confirmation; it appears in **My Enrollments**.
2. Try to enrol in the **same** activity again → blocked with a readable reason (not a raw error).
3. Reload the page → the enrolment is still there (it persisted).
4. Enrol in a second activity → both appear; the dashboard count increases.
5. Double-click the enrol button quickly → **only one** enrolment is created.

### S3.4 My enrolments

1. Statuses are consistent with the activities.
2. Any "withdraw"/"cancel" control works and persists, or is clearly disabled.
3. Empty state for a student with no enrolments is friendly.

**Findings:** *(table)*

---

## S4 — Student: certificate upload & management

**Tester:** ______  **Browser:** ______  **Date:** ______

This is the most important student flow. Use your **own** registered account where possible.

### S4.1 Upload — happy path

1. Portal → **Upload Certificate**. All required fields are marked.
2. Upload a real PDF with every field filled → success.
3. It appears in **My Certificates** with status **Pending**.
4. Reload → still there.

### S4.2 Upload — file validation

| #   | File                                           | Expected                                                            |
| --- | ---------------------------------------------- | ------------------------------------------------------------------- |
| a   | Valid PDF                                      | Accepted                                                            |
| b   | Valid JPG / PNG                                | Accepted                                                            |
| c   | `.txt` or `.docx`                              | Rejected, states allowed types                                      |
| d   | **A `.txt` renamed to `.pdf`**                 | **Rejected** — the server inspects content, not just the extension |
| e   | File larger than 5 MB                          | Rejected, states the limit                                          |
| f   | 0-byte file                                    | Rejected gracefully                                                 |
| g   | Filename with spaces/unicode (`my cert ग.pdf`) | Accepted, name shown correctly                                      |

### S4.3 Upload — field validation

1. Submit with **no file** → rejected.
2. Submit with a file but **no activity name / date / participation type** → rejected naming the missing field. It must never say "required fields are missing" while every visible field is filled.
3. Activity date in the **future** → report the behaviour (expected: rejected or flagged).
4. Issue date **earlier** than the activity date → report the behaviour.
5. Very long description (2000+ chars) → handled without breaking layout.
6. Press submit **twice** rapidly → only one certificate is created.

### S4.4 Viewing & downloading

1. Download your own certificate → the correct file opens.
2. Once S8 has rejected one of your certificates, confirm it is clearly marked **Rejected** **and shows the reason**.
3. Status filters/tabs return the right subsets.
4. Empty state for a student with no certificates is friendly.

### S4.5 Re-upload after rejection

1. For a certificate that S8 rejected, try to upload a replacement.
2. Report whether this is possible and whether it is obvious how — replacing a rejected certificate is a required behaviour, so if there is no way to do it, log it as **High**.

**Findings:** *(table)*

---

## S5 — Student: credits, marksheet & PDF export

**Tester:** ______  **Browser:** ______  **Date:** ______

### S5.1 Credits & progress

1. Credits shown match the dashboard for the **same** account — any disagreement between two screens is a finding.
2. Only **Approved** certificates contribute credits; Pending/Rejected must not.
3. The progress bar and percentage agree with the numbers beside them.
4. The target/graduation credit figure is displayed consistently everywhere it appears.
5. The category breakdown adds up to the total.

### S5.2 Extracurricular marksheet (on screen)

1. Rows correspond to real approved activities for that student.
2. Student name, roll number, course, semester and academic year are correct.
3. Totals at the bottom match the rows above.
4. For `vikram.singh` (the untouched control) the marksheet renders an empty state rather than a broken table.

### S5.3 PDF export — **priority for this round**

This was recently rebuilt; test it carefully.

1. Click **Download** / **Print**. In the browser's print dialog choose "Save as PDF".
2. **The PDF must contain the whole marksheet** — not a blank page, not only the first visible portion.
3. **Colours and backgrounds must be preserved** (header bars, status colours), matching what is on screen.
4. The app chrome is **absent**: no sidebar, no top navigation bar, no toast, no buttons.
5. If content spans multiple pages, no row is cut in half across a page break.
6. Table headers, the signature/verification block and the footer all appear.
7. Repeat on a **second browser** — report any difference.
8. Repeat for a student with **several** approved activities and for `vikram.singh`, who has **none**.
9. Report honestly whether the result looks like an official document you would hand to a student. Attach the PDF.

**Findings:** *(table — attach the generated PDFs)*

---

## S6 — Student: leaderboard & profile

**Tester:** ______  **Browser:** ______  **Date:** ______

### S6.1 Leaderboard

1. Ranked students appear with points; your own row is highlighted.
2. The order is genuinely descending by score.
3. Ties are handled sensibly (report how).
4. The year filter changes the data.
5. **Cross-check:** your dashboard credits vs your leaderboard points. If they disagree, report it with both numbers — the two use different time windows and this has been a real bug before.
6. Category champions are populated and plausible.
7. **Privacy:** the leaderboard must not expose emails, phone numbers or other personal data beyond name/course/score.

### S6.2 Profile

1. Your real details are shown.
2. Where there is no profile photo you see **your own initials**, not a stock photograph of a stranger and not someone else's initials.
3. Edit a field → save → reload → the change persisted.
4. Invalid inputs (letters in the phone field, malformed email) → rejected.
5. Change password with the **wrong** current password → rejected.
6. Change password correctly → old one stops working. **Restore `Pass@123` afterwards.**
7. Profile completeness (if shown) is consistent with what is actually filled in.

**Findings:** *(table)*

---

## S7 — Admin: access, profile & student management

**Tester:** ______  **Browser:** ______  **Date:** ______

### S7.1 Access

1. `admin` / `Pass@123` at `/admin-portal` signs in.
2. Wrong password → refused.
3. A newly created admin (ask the S10 tester, or the maintainer) is forced to change password on first login and **cannot skip it**.

### S7.2 Student management — batch scoping (**security-relevant**)

1. As `admin`, the student list contains **only IT2K24** students (5).
2. As `admin2`, only **IT2K25** students (3).
3. Search and filters work; clearing restores the list.
4. Open a student → real certificates and enrolments for **that** student.
5. **Tampering test:** as `admin2`, note the URL of an IT2K24 student from step 1 and open it directly. You must be refused. **If you can see another batch's student, stop and report it as Critical.**

### S7.3 Admin profile

1. Shows the signed-in admin's real name and batch — not a hardcoded name.
2. Edit and save → persists after reload.
3. Change password works. **Restore `Pass@123`.**

### S7.4 Dashboard & notifications

1. Dashboard stats are plausible and consistent with the student list (e.g. it must not claim 24 students when the list shows 5).
2. Open the notifications bell. Entries must describe **real** state — pending certificate counts, students needing attention — and must not mention infrastructure this project does not run (SSH logins, S3 backups, DB sync warnings). Report anything invented.
3. With nothing outstanding, an "all caught up" style message appears rather than an empty box.

**Findings:** *(table)*

---

## S8 — Admin: certificate verification (the credit loop)

**Tester:** ______  **Browser:** ______  **Date:** ______

This closes the product's core loop: student uploads → mentor verifies → credits awarded → leaderboard moves. Test it end to end.

### S8.1 The queue

1. As `admin`, open **Certificate Verification**. Pending certificates for IT2K24 are listed.
2. Each row shows student name, activity, date and status.
3. Status filters (All / Pending / Approved / Rejected) return the right subsets.
4. Search works.
5. As `admin2`, the queue shows only IT2K25 items.

### S8.2 Viewing the document

1. Select a certificate and click **Download** / **View**.
2. **The actual uploaded file downloads and opens.** A success toast with no file is a **High** finding.
3. Try it for several certificates.

### S8.3 Approve — must persist

1. Note a Pending certificate uploaded in S4 and the owning student.
2. Approve it → confirmation.
3. **Reload the page.** It must still be Approved. *(Reverting to Pending after reload is a Critical finding.)*
4. Log in as that student → the certificate shows Approved **and its credits now count** on their dashboard.
5. Check the leaderboard — that student's standing should reflect the new credits.

### S8.4 Reject — must persist and explain

1. Reject a Pending certificate.
2. A **reason must be required** — try submitting an empty reason; it should be refused.
3. Reload → still Rejected, reason retained.
4. Log in as that student → they can see the rejection **and the reason**.

### S8.5 Edge cases

1. Approve an already-approved certificate → sensible behaviour, no duplicate credit.
2. Two browser tabs: approve in one, then act on the same item in the other → no crash or corrupt state.
3. Double-click Approve → applied once.
4. **Tampering:** as `admin2`, try to approve an IT2K24 certificate via a manipulated URL/ID → must be refused. **Critical if it succeeds.**

> **Note:** you are working with certificates created during this round, so there is no baseline to restore. Do record what you approved and rejected in your report so other testers know what state the data is in.

**Findings:** *(table)*

---

## S9 — Admin: activity monitoring & batch analytics

**Tester:** ______  **Browser:** ______  **Date:** ______

These were newly wired to real data — check the numbers are real, not decorative.

### S9.1 Activity monitoring

1. The page loads with statistics, not spinners.
2. Numbers are plausible against the activities that exist and your batch's students. With no data yet, zeros are correct — confirm they read as zero rather than blank.
3. Insights/"students requiring attention" reference **real** students from your batch only.
4. Any filters work.
5. Send a reminder to a student (if offered) → clear success/failure feedback. Note whether an email actually arrives.
6. Any per-activity edit persists after reload.

### S9.2 Batch analytics

1. Overview loads with real batch names (IT2K24, IT2K25).
2. Open a batch → detail matches the summary row.
3. Numbers reconcile with the student list from S7 (same batch, same student count).
4. Export a report → the file downloads and opens; contents match what is on screen.
5. Charts/graphs render with axes and labels; no overlapping text.
6. Any batch note/override you save persists after reload.
7. As `admin2`, you see only your batch.

**Findings:** *(table)*

---

## S10 — Super admin: users, activities & tracks

**Tester:** ______  **Browser:** ______  **Date:** ______

### S10.1 Access & dashboard

1. `superadmin` signs in at `/super-admin-portal`; a plain `admin` is refused.
2. The four stat cards show real platform totals. At the start that is 8 students and 3 admins, with 0 activities; the activity count must increase after you create one.
3. Notifications describe real platform state, not invented security alerts.
4. "Recent system activity" reflects genuine actions — perform an action (e.g. create a track) and confirm it appears.

### S10.2 User management

1. The registry lists seeded students and admins with correct roles.
2. Create a new **Admin** → appears in the list; note the temporary password.
3. Hand it to the S7 tester (or verify yourself) that first login forces a password change.
4. Edit a user → persists.
5. Delete your test user → gone after reload.
6. Try to delete your **own** superadmin account → blocked.
7. Create a user with a duplicate ID/email → rejected clearly.
8. Create a user with blank fields → rejected.

### S10.3 Activity management

1. Activity list matches the catalogue students see.
2. **Create** an activity → **reload** → it persists (this used to be local-only and vanish).
3. It now appears in a student's Browse Activities.
4. Edit it → persists. Delete it → gone.
5. Invalid input (blank name, negative credits, end date before start date) → rejected.
6. The track dropdown is populated.

### S10.4 Track management

1. Track list loads with activity counts (the page must not 404 — this was recently repaired).
2. Create a track → persists after reload.
3. Duplicate name → rejected as a conflict.
4. Edit description/status → persists.
5. Delete a track that still has activities → report exactly what happens.
6. Name over 100 characters → rejected.

**Findings:** *(table)*

---

## S11 — Super admin: announcements, settings & reports

**Tester:** ______  **Browser:** ______  **Date:** ______

### S11.1 Announcements

1. The announcement list starts empty with a friendly empty state.
2. Create a **draft** → persists.
3. **Publish** it → status changes and survives reload.
4. Edit and delete → persist.
5. Expiry date before publish date → rejected.
6. Empty title/body → rejected.
7. Check whether a published announcement is actually visible to a student — report if not.

### S11.2 System settings

1. Settings load, grouped by category (26 entries — settings were intentionally kept).
2. Change a value → save → reload → persisted.
3. Invalid values (letters where a number belongs, negatives) → rejected.
4. **If a "graduation target credits" setting exists**, change it and confirm the student Credits & Progress screen reflects the new target. Restore it afterwards.
5. Note which settings visibly change app behaviour and which appear display-only — useful information either way.

### S11.3 Reports centre

1. Summary, templates and the institutional overview load with real figures.
2. Generate a report → it completes and appears in the list.
3. Download it → the file opens and its contents match the filters you chose.
4. Apply filters (course, semester, date range) → results actually change.
5. A filter combination matching nothing → an empty report or clear message, not a crash.
6. Schedule a report → it is saved and listed. Delete it.
7. The audit log shows your actions from this session.
8. **Important:** generate a report, then wait for the API to sleep (~15 min idle) and try downloading the same report again. Report whether it still works — generated files may not survive a restart.

**Findings:** *(table)*

---

## S12 — Cross-cutting: security, responsive & resilience

**Tester:** ______  **Browser:** ______  **Date:** ______

Assign to someone who has already completed another suite.

### S12.1 Access control

1. Logged out, open each directly: `/portal`, `/admin-portal/dashboard`, `/super-admin-portal/dashboard`. No real data may render.
2. As a student, open an admin URL → refused.
3. As `admin`, open a super-admin URL → refused.
4. Log in as a student, then in DevTools → Application → Local Storage, note the token exists. Log out → confirm it is cleared.
5. Confirm no page displays another user's personal data anywhere.

### S12.2 Input safety

Try these in free-text fields (activity name, description, profile name, notes, search):

1. `<script>alert('x')</script>` → must be displayed as text, never executed as a popup.
2. `'; DROP TABLE students;--` → treated as ordinary text.
3. A 5000-character string → handled without breaking the page.
4. Emoji and non-Latin script → stored and displayed correctly.
5. Leading/trailing spaces → trimmed or handled consistently.

### S12.3 Responsive layout

Check **375px (phone)**, **768px (tablet)** and **1440px (desktop)** across the landing page, login, student portal (every tab), admin portal and super-admin portal:

1. No horizontal scrolling of the page body.
2. Navigation collapses to a working menu on small screens.
3. Wide tables scroll inside their own container rather than stretching the page.
4. Modals fit on screen and can always be closed.
5. Buttons are large enough to tap and are not overlapped.
6. Text never overflows its container.

### S12.4 Resilience

1. **Offline:** DevTools → Network → Offline, then click something that loads data. Expect a readable error, not a silent hang or a raw stack trace.
2. **Slow network:** set throttling to "Slow 3G" and load the portal → loading states appear; nothing renders as `undefined`/`NaN`.
3. **Cold start:** after ~15 minutes idle, the first login is slow but must **succeed** and show a loading indicator rather than an unexplained error.
4. **Back/forward:** navigate several tabs, then use browser Back and Forward repeatedly → no broken state.
5. **Double submit:** on any form, double-click submit → no duplicate record.
6. **Session expiry:** leave a portal tab open for a long period, then act → either it still works or you are asked to sign in again. It must not half-work or show a raw error.

### S12.5 Consistency sweep

1. The same student's credit total is identical on the dashboard, Credits & Progress and the marksheet.
2. Certificate counts agree between the student's view and the admin's view.
3. Dates use a consistent format across the app.
4. Terminology is consistent (e.g. "Enrolment" vs "Registration"; "Mentor" vs "Admin") — list inconsistencies.
5. Every number a page shows can be traced to something real. Flag any figure you cannot explain.

**Findings:** *(table)*

---

## 5. Sign-off

| Suite | Tester | Browser(s) | Date | Result | Critical/High open |
| ----- | ------ | ---------- | ---- | ------ | ------------------ |
| S1    |        |            |      |        |                    |
| S2    |        |            |      |        |                    |
| S3    |        |            |      |        |                    |
| S4    |        |            |      |        |                    |
| S5    |        |            |      |        |                    |
| S6    |        |            |      |        |                    |
| S7    |        |            |      |        |                    |
| S8    |        |            |      |        |                    |
| S9    |        |            |      |        |                    |
| S10   |        |            |      |        |                    |
| S11   |        |            |      |        |                    |
| S12   |        |            |      |        |                    |

**Definition of done for a suite:** every case executed on at least one browser, every finding logged with steps and a screenshot, and all Critical/High items reported to the maintainer directly.
