-- Clear all content so contributors start from an empty canvas, while keeping
-- the accounts they need in order to sign in and test.
--
--   node scripts/db.mjs --file scripts/sql/reset-content-keep-accounts.sql
--
-- KEPT
--   students, admins   the login credentials themselves are under test
--   system_settings    platform configuration, not test data; the Settings
--                      screen needs rows to be testable at all
--
-- CLEARED
--   certificates, enrollments, activities, tracks, announcements,
--   admin_notes, batch_overrides, and all report history + audit logs, plus any
--   pending OTPs.
--
-- This deliberately leaves the platform with no tracks and no activities, so the
-- natural test order becomes the real product flow:
--   super admin creates a track -> creates an activity -> student enrols ->
--   student uploads a certificate -> admin verifies it -> credits appear.

truncate table
  certificates,
  enrollments,
  admin_notes,
  activities,
  tracks,
  announcements,
  batch_overrides,
  generated_reports,
  scheduled_reports,
  report_audit_logs,
  otps
restart identity cascade;
