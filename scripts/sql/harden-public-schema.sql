-- Lock down the public schema against Supabase's automatically exposed REST API.
--
-- Supabase serves every table in the `public` schema over PostgREST at
-- https://<ref>.supabase.co/rest/v1/, authorised by the `anon` key — which is
-- public by design and ships inside client applications. iSPARC does not use
-- Supabase Auth or PostgREST at all: the Go API owns every access-control
-- decision and connects as the table owner. Any privilege held by `anon` or
-- `authenticated` is therefore pure attack surface, and without RLS it exposed
-- student password hashes and live OTP codes to anyone with that public key.
--
-- Enabling RLS with no policies denies all non-owner roles. The API is
-- unaffected because the table owner bypasses RLS unless it is FORCEd — which
-- this script deliberately does not do.
--
-- Re-run after a schema change adds tables:
--   node scripts/db.mjs --file scripts/sql/harden-public-schema.sql

do $$
declare rec record;
begin
  for rec in
    select c.relname
    from pg_class c
    join pg_namespace n on n.oid = c.relnamespace
    where n.nspname = 'public' and c.relkind = 'r'
  loop
    execute format('alter table public.%I enable row level security', rec.relname);
    execute format('revoke all on public.%I from anon, authenticated', rec.relname);
  end loop;
end $$;

-- Tables created later (for example by GORM AutoMigrate after a model change)
-- must not inherit access either.
alter default privileges in schema public revoke all on tables from anon, authenticated;
alter default privileges in schema public revoke all on sequences from anon, authenticated;
alter default privileges in schema public revoke all on functions from anon, authenticated;

-- Nothing in this project reaches the database through PostgREST, so the
-- anonymous roles have no reason to see the schema at all. Re-grant this if you
-- ever decide to use the Supabase client libraries directly from the browser.
revoke usage on schema public from anon, authenticated;
