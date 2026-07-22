-- Remove rows left behind by scripts/api-regression.mjs when it is run against
-- a persistent (cloud) database. The suite registers throwaway students,
-- creates temporary admins, and adds sample tracks/announcements; against the
-- local disposable Postgres this does not matter, but against Supabase these
-- accumulate and pollute the demo data. Safe to run any time.

-- Regression students use roll numbers starting RT and reg.test / plus-tagged
-- emails. Remove their dependent rows first.
delete from certificates where student_roll_no in (select roll_no from students where roll_no like 'RT%');
delete from enrollments  where student_roll_no in (select roll_no from students where roll_no like 'RT%');
delete from students where roll_no like 'RT%';
delete from otps where email like 'reg.test%' or email like '%+ispark%' or email like '%@example.com';

-- Temporary admins from the "create platform user" test.
delete from admins where admin_id like 'tmp%';

-- Sample tracks / announcements created by the super-admin CRUD tests.
delete from tracks where name like 'Regression Track%';
delete from announcements where title like 'Regression Announcement%';
