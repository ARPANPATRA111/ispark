// Run SQL against the project's database without needing Docker or psql.
//
//   node scripts/db.mjs "select count(*) from students"          # cloud (api/.env.supabase)
//   node scripts/db.mjs --local "select count(*) from students"  # local  (api/.env)
//   node scripts/db.mjs --file scripts/sql/rls.sql               # run a .sql file
//
// The connection is read from the gitignored env files, so no credentials live
// in this script or anywhere in the repository.
import pg from 'pg';
import fs from 'node:fs';
import path from 'node:path';
import { fileURLToPath } from 'node:url';

const repoRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), '..');

function readEnvFile(file) {
	const values = {};
	if (!fs.existsSync(file)) return values;
	for (const line of fs.readFileSync(file, 'utf8').split(/\r?\n/)) {
		const match = /^([A-Z0-9_]+)=(.*)$/.exec(line.trim());
		if (match) values[match[1]] = match[2].trim();
	}
	return values;
}

function connectionString({ local }) {
	if (process.env.DATABASE_URL) return process.env.DATABASE_URL.trim();

	const file = path.join(repoRoot, 'api', local ? '.env' : '.env.supabase');
	const env = readEnvFile(file);
	if (env.DATABASE_URL) return env.DATABASE_URL;

	const { DB_USER, DB_PASSWORD, DB_HOST, DB_PORT, DB_NAME, DB_SSLMODE } = env;
	if (!DB_HOST) throw new Error(`No database configuration found in ${file}`);

	const sslmode = DB_SSLMODE ? `?sslmode=${DB_SSLMODE}` : '';
	return `postgresql://${encodeURIComponent(DB_USER)}:${encodeURIComponent(DB_PASSWORD)}@${DB_HOST}:${DB_PORT || 5432}/${DB_NAME}${sslmode}`;
}

const args = process.argv.slice(2);
const local = args.includes('--local');
const fileFlag = args.indexOf('--file');
const sql =
	fileFlag !== -1
		? fs.readFileSync(path.resolve(repoRoot, args[fileFlag + 1]), 'utf8')
		: args.find((a) => !a.startsWith('--'));

if (!sql) {
	console.error('Provide a SQL string or --file <path>.');
	process.exit(1);
}

const rawConn = connectionString({ local });

// Supabase's pooler presents a chain Node does not trust out of the box. The
// API connects with sslmode=require, which encrypts but does not verify the
// certificate, so this tool matches that posture rather than being stricter
// than the application it inspects. sslmode is stripped from the URL because
// node-postgres would otherwise derive its own (verifying) TLS settings from it.
const usesTls = /sslmode=(require|verify)/.test(rawConn);
const conn = rawConn.replace(/[?&]sslmode=[^&]*/g, '');
const ssl = usesTls ? { rejectUnauthorized: false } : undefined;

const client = new pg.Client({ connectionString: conn, ssl });
await client.connect();
try {
	const result = await client.query(sql);
	for (const res of Array.isArray(result) ? result : [result]) {
		if (res.rows?.length) console.table(res.rows);
		else console.log(`${res.command ?? 'OK'} (${res.rowCount ?? 0} rows)`);
	}
} finally {
	await client.end();
}
