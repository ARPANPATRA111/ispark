-- Keep the free-tier database and API awake.
--
-- Two independent mechanisms, so neither is a single point of failure:
--
--   1. An external uptime monitor calls GET /health/db on the API. That one
--      request resets the API host's idle timer AND runs a real query against
--      this database, which is the activity signal that prevents pausing.
--
--   2. The pg_cron jobs below, which work even if the monitor lapses:
--        * a heartbeat row written and later deleted, so the database sees
--          genuine write activity on a fixed schedule;
--        * an HTTP call to the API's health endpoint via pg_net, so the
--          database keeps the API awake without any third-party service.
--
-- Every pg_cron run also inserts into cron.job_run_details, so frequent jobs
-- generate write activity in their own right.
--
-- Apply with:
--   node scripts/db.mjs --file scripts/sql/keepalive.sql
--
-- Schedules are in the database's timezone (UTC on Supabase). 10:00 UTC is
-- 15:30 IST; subtract 5h30m from a desired IST time to get the UTC value.

create extension if not exists pg_cron;
create extension if not exists pg_net;

-- A dedicated schema: unlike `public`, it is not exposed through PostgREST, so
-- this table can never be reached from the internet regardless of RLS.
create schema if not exists ops;
revoke all on schema ops from anon, authenticated;

create table if not exists ops.keepalive (
  id          bigserial primary key,
  note        text        not null,
  created_at  timestamptz not null default now()
);

alter table ops.keepalive enable row level security;
revoke all on ops.keepalive from anon, authenticated;

-- Re-running this file should not stack duplicate jobs.
select cron.unschedule(jobname)
from cron.job
where jobname in ('keepalive-write', 'keepalive-delete', 'keepalive-ping-api');

-- Write a heartbeat row at 10:00 on days 1, 6, 11, 16, 21 and 26.
select cron.schedule(
  'keepalive-write',
  '0 10 1-31/5 * *',
  $$insert into ops.keepalive (note) values ('heartbeat')$$
);

-- Remove it at 08:00 the following day (days 2, 7, 12, 17, 22 and 27), which
-- is a second write-type operation and keeps the table from growing.
select cron.schedule(
  'keepalive-delete',
  '0 8 2-31/5 * *',
  $$delete from ops.keepalive where created_at < now() - interval '12 hours'$$
);

-- Call the API every 10 minutes so it never reaches the 15-minute idle
-- threshold. /health/db also queries this database, so the request keeps both
-- halves of the stack alive.
select cron.schedule(
  'keepalive-ping-api',
  '*/10 * * * *',
  $$select net.http_get(
      url := 'https://ispark-api.onrender.com/health/db',
      timeout_milliseconds := 20000
    )$$
);
