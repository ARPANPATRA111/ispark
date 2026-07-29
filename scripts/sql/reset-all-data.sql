-- Wipe every application row so a testing round starts from a known baseline.
--
-- DESTRUCTIVE: this deletes all students, admins, certificates, enrolments,
-- activities, tracks, announcements, settings and report history. Run it only
-- against a testing database, never production.
--
--   node scripts/db.mjs --file scripts/sql/reset-all-data.sql
--
-- Afterwards boot the API once with SEED_DEV_DATA=true to recreate the demo
-- dataset, then set it back to false.
--
-- TRUNCATE ... RESTART IDENTITY resets the id sequences too, so the seeded rows
-- get predictable ids (certificate 1..10) that the manual test plan can refer
-- to. CASCADE covers the foreign keys between activities and tracks.

truncate table
  certificates,
  enrollments,
  otps,
  admin_notes,
  students,
  admins,
  activities,
  tracks,
  announcements,
  system_settings,
  batch_overrides,
  generated_reports,
  scheduled_reports,
  report_audit_logs
restart identity cascade;
