// iSPARC API regression suite — run with: node scripts/api-regression.mjs
// Requires: API running (default http://localhost:8080, override with API_BASE_URL)
// and the docker container ispark-db-1 running (for OTP lookup in register/reset flows).
import pg from 'pg';

const API = process.env.API_BASE_URL || 'http://localhost:8080';
const PASS_WORD = 'Pass@123';
let passed = 0, failed = 0;
const failures = [];

function check(name, cond, detail = '') {
  if (cond) { passed++; console.log(`PASS  ${name}`); }
  else { failed++; failures.push(`${name} ${detail}`); console.log(`FAIL  ${name}  ${detail}`); }
}

async function j(method, path, { token, body, form, cookies } = {}) {
  const headers = {};
  if (token) headers.Authorization = `Bearer ${token}`;
  if (cookies) headers.Cookie = cookies;
  let payload;
  if (form) { payload = form; }
  else if (body !== undefined) { headers['Content-Type'] = 'application/json'; payload = JSON.stringify(body); }
  const res = await fetch(`${API}${path}`, { method, headers, body: payload });
  const ct = res.headers.get('content-type') || '';
  const setCookie = res.headers.getSetCookie ? res.headers.getSetCookie() : [];
  let data = null, raw = null;
  if (ct.includes('application/json')) data = await res.json();
  else raw = await res.arrayBuffer();
  return { status: res.status, data, raw, setCookie, headers: res.headers };
}

// Address used for the throwaway registration/reset accounts.
//
// Against a deployment with live SMTP configured, mail to example.com hard
// bounces and damages the sender reputation the real OTPs depend on. Set
// TEST_EMAIL_BASE to an inbox you own and the suite uses plus-addressing, so
// every test message is deliverable and lands in that one inbox.
function makeTestEmail() {
  const tag = `ispark${Date.now()}`;
  const base = process.env.TEST_EMAIL_BASE;
  if (base && base.includes('@')) {
    const [local, domain] = base.split('@');
    return `${local}+${tag}@${domain}`;
  }
  return `reg.test.${tag}@example.com`;
}

async function solveCaptcha() {
  const { data } = await j('GET', '/api/auth/captcha');
  const m = data.question.match(/(\d+)\s*\+\s*(\d+)/);
  return { captcha_id: data.captcha_id, captcha_answer: String(Number(m[1]) + Number(m[2])) };
}

async function studentLogin(email, password = PASS_WORD) {
  const cap = await solveCaptcha();
  return j('POST', '/api/auth/login', { body: { email_id: email, password, ...cap } });
}

// Reads the OTP straight from the database. OTP_DB_URL points at whichever
// database the API under test is using (e.g. the Supabase session pooler);
// it defaults to the local docker-compose Postgres.
async function otpFromDb(email, purpose) {
  const raw = process.env.OTP_DB_URL || 'postgresql://postgres:user@localhost:5432/isparc';
  const usesTls = /sslmode=(require|verify)/.test(raw);
  const client = new pg.Client({
    connectionString: raw.replace(/[?&]sslmode=[^&]*/g, ''),
    ssl: usesTls ? { rejectUnauthorized: false } : undefined
  });

  await client.connect();
  try {
    const res = await client.query(
      'select code from otps where email = $1 and purpose = $2 order by created_at desc limit 1',
      [email, purpose]
    );
    return res.rows[0]?.code ?? '';
  } finally {
    await client.end();
  }
}

// Minimal valid PDF for upload content-sniffing
const pdfBytes = Buffer.from(
  '%PDF-1.4\n1 0 obj<</Type/Catalog/Pages 2 0 R>>endobj\n2 0 obj<</Type/Pages/Kids[3 0 R]/Count 1>>endobj\n' +
  '3 0 obj<</Type/Page/Parent 2 0 R/MediaBox[0 0 200 200]>>endobj\nxref\n0 4\ntrailer<</Size 4/Root 1 0 R>>\n%%EOF\n'
);

const run = async () => {
console.log('=== 1. AUTH — student =========================================');

// 1.1 captcha
{
  const { status, data } = await j('GET', '/api/auth/captcha');
  check('1.1 GET /auth/captcha returns question + id', status === 200 && !!data.captcha_id && /\d+ \+ \d+/.test(data.question));
}

// 1.2 login happy path (seeded student)
let studentToken, studentCookies;
{
  const r = await studentLogin('rahul.sharma@iips.edu');
  studentToken = r.data?.access_token;
  studentCookies = (r.setCookie || []).map(c => c.split(';')[0]).join('; ');
  check('1.2 student login (rahul.sharma) returns access_token', r.status === 200 && !!studentToken, `status=${r.status} body=${JSON.stringify(r.data)}`);
  check('1.3 login sets refresh_token cookie', (r.setCookie || []).some(c => c.startsWith('refresh_token=')));
}

// 1.4 wrong password
{
  const r = await studentLogin('rahul.sharma@iips.edu', 'WrongPass@1');
  check('1.4 wrong password rejected 401', r.status === 401);
}

// 1.5 missing captcha
{
  const r = await j('POST', '/api/auth/login', { body: { email_id: 'rahul.sharma@iips.edu', password: PASS_WORD } });
  check('1.5 login without captcha rejected 400', r.status === 400);
}

// 1.6 wrong captcha answer
{
  const cap = await solveCaptcha();
  const r = await j('POST', '/api/auth/login', { body: { email_id: 'rahul.sharma@iips.edu', password: PASS_WORD, captcha_id: cap.captcha_id, captcha_answer: '99999' } });
  check('1.6 wrong captcha answer rejected 400', r.status === 400);
}

// 1.7 register → OTP → verify
const regEmail = makeTestEmail();
let regToken;
{
  const r = await j('POST', '/api/auth/register', { body: {
    name: 'Regression Tester', roll_no: `RT${Date.now() % 100000}`, course_name: 'MCA', semester: 2,
    contact_no: '9999999999', email_id: regEmail, enrollment_no: `EN-RT${Date.now() % 100000}`,
    password: 'Student@123', confirm_password: 'Student@123'
  }});
  check('1.7 register new student accepted', r.status === 200 || r.status === 201, `status=${r.status} ${JSON.stringify(r.data)}`);
  const otp = await otpFromDb(regEmail, 'register');
  check('1.8 OTP row created in DB', /^\d{6}$/.test(otp), `otp=${otp}`);
  const v = await j('POST', '/api/auth/verify-otp', { body: { email: regEmail, code: otp } });
  regToken = v.data?.access_token;
  check('1.9 verify-otp returns access_token', v.status === 200 && !!regToken, `status=${v.status} ${JSON.stringify(v.data)}`);
}

// 1.10 duplicate register
{
  const r = await j('POST', '/api/auth/register', { body: {
    name: 'Dup', roll_no: 'IT2K24-01', course_name: 'MCA', semester: 2, contact_no: '9999999999',
    email_id: 'rahul.sharma@iips.edu', enrollment_no: 'DUP-1', password: 'Student@123', confirm_password: 'Student@123'
  }});
  check('1.10 duplicate email register rejected', r.status >= 400, `status=${r.status}`);
}

// 1.11 profile with/without token
{
  const a = await j('GET', '/api/auth/profile', { token: studentToken });
  check('1.11 GET /auth/profile with token', a.status === 200 && (a.data?.student?.email_id === 'rahul.sharma@iips.edu' || a.data?.email_id === 'rahul.sharma@iips.edu' || JSON.stringify(a.data).includes('rahul.sharma')), `status=${a.status}`);
  const b = await j('GET', '/api/auth/profile');
  check('1.12 GET /auth/profile without token → 401', b.status === 401);
}

// 1.13 refresh token
{
  const r = await j('POST', '/api/auth/refresh', { cookies: studentCookies });
  check('1.13 POST /auth/refresh with cookie returns new access token', r.status === 200 && !!r.data?.access_token, `status=${r.status} ${JSON.stringify(r.data)}`);
}

// 1.14 forgot/reset password (on freshly registered account)
{
  const f = await j('POST', '/api/auth/forgot-password', { body: { email_id: regEmail } });
  check('1.14 forgot-password accepted', f.status === 200, `status=${f.status} ${JSON.stringify(f.data)}`);
  const otp = await otpFromDb(regEmail, 'forgot_password');
  check('1.15 reset OTP created', /^\d{6}$/.test(otp), `otp=${otp}`);
  const rs = await j('POST', '/api/auth/reset-password', { body: { email_id: regEmail, otp_code: otp, password: 'NewPass@123', confirm_password: 'NewPass@123' } });
  check('1.16 reset-password succeeds', rs.status === 200, `status=${rs.status} ${JSON.stringify(rs.data)}`);
  const l = await studentLogin(regEmail, 'NewPass@123');
  check('1.17 login works with new password', l.status === 200 && !!l.data?.access_token, `status=${l.status}`);
  regToken = l.data?.access_token || regToken;
}

console.log('=== 2. STUDENT MODULE =========================================');

// 2.1 dashboard stats
{
  const r = await j('GET', '/api/student/dashboard/stats', { token: studentToken });
  check('2.1 dashboard stats 200 with credit fields', r.status === 200 && JSON.stringify(r.data).includes('credits'), `status=${r.status}`);
}

// 2.2 activities catalogue
let activities = [];
{
  const r = await j('GET', '/api/student/activities', { token: studentToken });
  activities = r.data?.activities || r.data || [];
  check('2.2 activities catalogue returns 7 seeded activities', r.status === 200 && (activities.length ?? 0) >= 7, `status=${r.status} count=${activities.length}`);
}

// 2.3 enroll (fresh student, first activity)
{
  const id = activities[0]?.id ?? activities[0]?.ID;
  const r = await j('POST', `/api/student/activities/${id}/enroll`, { token: regToken });
  check('2.3 fresh student can enroll', r.status === 200 || r.status === 201, `status=${r.status} ${JSON.stringify(r.data)}`);
  const dup = await j('POST', `/api/student/activities/${id}/enroll`, { token: regToken });
  check('2.4 duplicate enrollment rejected', dup.status >= 400, `status=${dup.status}`);
  const e = await j('GET', '/api/student/enrollments', { token: regToken });
  const list = e.data?.enrollments || e.data || [];
  check('2.5 enrollments list shows the new enrollment', e.status === 200 && list.length >= 1, `status=${e.status} count=${list.length}`);
}

// 2.6 certificate upload (fresh student)
let certId;
{
  const form = new FormData();
  form.set('activity_name', 'Regression Workshop');
  form.set('activity_category', 'TECHNICAL');
  form.set('activity_date', '2026-07-01');
  form.set('organizer_name', 'QA Cell');
  form.set('event_level', 'College');
  form.set('cert_number', `CERT-REG-${Date.now()}`);
  form.set('issue_date', '2026-07-02');
  form.set('participation_type', 'Participant');
  form.set('description', 'Uploaded by automated regression suite');
  form.set('certificate_file', new Blob([pdfBytes], { type: 'application/pdf' }), 'regression.pdf');
  const r = await j('POST', '/api/student/certificates', { token: regToken, form });
  check('2.6 certificate upload accepted', r.status === 200 || r.status === 201, `status=${r.status} ${JSON.stringify(r.data)}`);
  const list = await j('GET', '/api/student/certificates', { token: regToken });
  const certs = list.data?.certificates || list.data || [];
  certId = certs[0]?.id ?? certs[0]?.ID;
  check('2.7 uploaded certificate appears in list (Pending)', list.status === 200 && certs.length >= 1, `count=${certs.length}`);
}

// 2.8 non-PDF upload rejected
{
  const form = new FormData();
  form.set('activity_name', 'Bad Upload');
  form.set('activity_category', 'TECHNICAL');
  form.set('activity_date', '2026-07-01');
  form.set('organizer_name', 'QA Cell');
  form.set('event_level', 'College');
  form.set('cert_number', `CERT-BAD-${Date.now()}`);
  form.set('issue_date', '2026-07-02');
  form.set('participation_type', 'Participant');
  form.set('certificate_file', new Blob([Buffer.from('MZ executable junk')], { type: 'application/pdf' }), 'evil.pdf');
  const r = await j('POST', '/api/student/certificates', { token: regToken, form });
  check('2.8 fake-PDF (bad magic bytes) rejected', r.status >= 400, `status=${r.status}`);
}

// 2.9 download own certificate
{
  const r = await j('GET', `/api/student/certificates/${certId}/file`, { token: regToken });
  const head = r.raw ? Buffer.from(r.raw.slice(0, 4)).toString() : '';
  check('2.9 own certificate downloads as PDF', r.status === 200 && head === '%PDF', `status=${r.status} head=${head}`);
}

// 2.10 cross-student certificate access blocked
{
  const other = await j('GET', '/api/student/certificates', { token: studentToken });
  const otherCerts = other.data?.certificates || other.data || [];
  const otherId = otherCerts[0]?.id ?? otherCerts[0]?.ID;
  if (otherId) {
    const r = await j('GET', `/api/student/certificates/${otherId}/file`, { token: regToken });
    check('2.10 downloading another student\'s certificate blocked', r.status === 403 || r.status === 404, `status=${r.status}`);
  } else {
    check('2.10 downloading another student\'s certificate blocked', false, 'no cert found for rahul to test against');
  }
}

// 2.11 leaderboard + champions
{
  const r = await j('GET', '/api/student/leaderboard', { token: studentToken });
  const rows = r.data?.leaderboard || r.data || [];
  check('2.11 leaderboard returns ranked rows', r.status === 200 && rows.length >= 1, `status=${r.status} count=${rows.length}`);
  const c = await j('GET', '/api/student/leaderboard/champions', { token: studentToken });
  check('2.12 category champions 200', c.status === 200, `status=${c.status}`);
}

// 2.13 marksheet
{
  const r = await j('GET', '/api/student/marksheet', { token: studentToken });
  check('2.13 marksheet 200', r.status === 200, `status=${r.status}`);
}

// 2.14 profile update + change password (fresh student)
{
  const r = await j('PUT', '/api/student/profile', { token: regToken, body: { contact_no: '8888888888' } });
  check('2.14 profile update 200', r.status === 200, `status=${r.status} ${JSON.stringify(r.data)}`);
  const cp = await j('POST', '/api/student/change-password', { token: regToken, body: { current_password: 'NewPass@123', new_password: 'Final@123', confirm_password: 'Final@123' } });
  check('2.15 change-password 200', cp.status === 200, `status=${cp.status} ${JSON.stringify(cp.data)}`);
  const l = await studentLogin(regEmail, 'Final@123');
  check('2.16 login with changed password works', l.status === 200, `status=${l.status}`);
}

console.log('=== 3. RBAC BOUNDARIES ========================================');

{
  const r = await j('GET', '/api/admin/students', { token: studentToken });
  check('3.1 student token on admin route → 403', r.status === 403, `status=${r.status}`);
  const r2 = await j('GET', '/api/student/dashboard/stats');
  check('3.2 no token on student route → 401', r2.status === 401, `status=${r2.status}`);
}

console.log('=== 4. ADMIN MODULE ===========================================');

let adminToken;
{
  const r = await j('POST', '/api/admin/auth/login', { body: { admin_id: 'admin', password: PASS_WORD } });
  adminToken = r.data?.access_token;
  check('4.1 admin login returns token', r.status === 200 && !!adminToken, `status=${r.status} ${JSON.stringify(r.data)}`);
  const bad = await j('POST', '/api/admin/auth/login', { body: { admin_id: 'admin', password: 'nope' } });
  check('4.2 admin wrong password → 401', bad.status === 401, `status=${bad.status}`);
}

{
  const r = await j('GET', '/api/admin/profile', { token: adminToken });
  check('4.3 admin profile 200', r.status === 200, `status=${r.status}`);
}

let batchRoll, otherBatchRoll;
{
  const r = await j('GET', '/api/admin/students', { token: adminToken });
  const students = r.data?.students || r.data || [];
  batchRoll = students[0]?.roll_no ?? students[0]?.RollNo;
  check('4.4 admin sees only own batch (IT2K24 → 5 seeded)', r.status === 200 && students.length >= 5, `status=${r.status} count=${students.length}`);
  const all = students.map(s => s.batch ?? s.Batch).filter(Boolean);
  check('4.5 all listed students share admin batch', all.length === 0 || all.every(b => b === all[0]), `batches=${[...new Set(all)]}`);
}

{
  const r = await j('GET', `/api/admin/students/${batchRoll}`, { token: adminToken });
  check('4.6 student detail for own batch 200', r.status === 200, `status=${r.status}`);
}

// admin2 (IT2K25) cross-batch check
{
  const r2 = await j('POST', '/api/admin/auth/login', { body: { admin_id: 'admin2', password: PASS_WORD } });
  const admin2Token = r2.data?.access_token;
  const cross = await j('GET', `/api/admin/students/${batchRoll}`, { token: admin2Token });
  check('4.7 admin2 cannot read IT2K24 student detail', cross.status === 403 || cross.status === 404, `status=${cross.status}`);
  const list2 = await j('GET', '/api/admin/students', { token: admin2Token });
  const students2 = list2.data?.students || list2.data || [];
  otherBatchRoll = students2[0]?.roll_no;
  check('4.8 admin2 sees IT2K25 batch (3 seeded)', list2.status === 200 && students2.length >= 3, `count=${students2.length}`);
}

// 4.9 admin blocked from platform routes
{
  const r = await j('GET', '/api/admin/platform/stats', { token: adminToken });
  check('4.9 admin token on superadmin platform route → 403', r.status === 403, `status=${r.status}`);
}

console.log('=== 5. SUPER ADMIN MODULE =====================================');

let saToken;
{
  const r = await j('POST', '/api/admin/auth/login', { body: { admin_id: 'superadmin', password: PASS_WORD } });
  saToken = r.data?.access_token;
  check('5.1 superadmin login returns token', r.status === 200 && !!saToken, `status=${r.status}`);
}

{
  const r = await j('GET', '/api/admin/platform/stats', { token: saToken });
  check('5.2 platform stats 200', r.status === 200 && JSON.stringify(r.data).match(/student|user/i), `status=${r.status}`);
  const u = await j('GET', '/api/admin/platform/users', { token: saToken });
  check('5.3 platform users list 200', u.status === 200, `status=${u.status}`);
}

// 5.4 user create/delete (delete key is the business id: admin_id / roll_no)
{
  const tmpId = `tmp${Date.now() % 100000}`;
  const c = await j('POST', '/api/admin/platform/users', { token: saToken, body: { name: 'Temp Mentor', role: 'Admin', id: tmpId, email: `tmp${Date.now()}@iips.edu`, dept: 'IT2K24' } });
  check('5.4 create platform user (Admin)', c.status === 200 || c.status === 201, `status=${c.status} ${JSON.stringify(c.data)}`);
  const d = await j('DELETE', `/api/admin/platform/users/${tmpId}`, { token: saToken });
  check('5.5 delete platform user', d.status === 200, `status=${d.status} ${JSON.stringify(d.data)}`);
}

// 5.6 settings
{
  const g = await j('GET', '/api/admin/platform/settings', { token: saToken });
  check('5.6 get settings 200', g.status === 200, `status=${g.status}`);
  // Response shape: { settings: { "<Category>": [ {key, value, ...} ] } }
  const grouped = g.data?.settings || {};
  const first = Object.values(grouped).flat()[0] ?? null;
  if (first?.key) {
    const upd = await j('PUT', `/api/admin/platform/settings/${first.key}`, { token: saToken, body: { value: first.value } });
    check('5.7 update single setting 200', upd.status === 200, `status=${upd.status} ${JSON.stringify(upd.data)}`);
  } else {
    check('5.7 update single setting 200', false, `no settings rows: ${JSON.stringify(g.data).slice(0, 120)}`);
  }
}

// 5.8 tracks CRUD
{
  const stats = await j('GET', '/api/admin/platform/tracks/stats', { token: saToken });
  check('5.8 track stats 200', stats.status === 200, `status=${stats.status}`);
  const list = await j('GET', '/api/admin/platform/tracks', { token: saToken });
  check('5.9 tracks list 200', list.status === 200, `status=${list.status}`);
  const c = await j('POST', '/api/admin/platform/tracks', { token: saToken, body: { name: `Regression Track ${Date.now()}`, description: 'temp', status: 'Active' } });
  const tid = (c.data?.track ?? c.data)?.id;
  check('5.10 create track', c.status === 200 || c.status === 201, `status=${c.status} ${JSON.stringify(c.data)}`);
  if (tid) {
    const u = await j('PUT', `/api/admin/platform/tracks/${tid}`, { token: saToken, body: { description: 'updated by regression' } });
    check('5.11 update track', u.status === 200, `status=${u.status}`);
    const d = await j('DELETE', `/api/admin/platform/tracks/${tid}`, { token: saToken });
    check('5.12 delete track', d.status === 200, `status=${d.status}`);
  } else {
    check('5.11 update track', false, 'no track id'); check('5.12 delete track', false, 'no track id');
  }
}

// 5.13 announcements CRUD + publish
{
  const stats = await j('GET', '/api/admin/platform/announcements/stats', { token: saToken });
  check('5.13 announcement stats 200', stats.status === 200, `status=${stats.status}`);
  const list = await j('GET', '/api/admin/platform/announcements', { token: saToken });
  check('5.14 announcements list 200 (seeded rows)', list.status === 200, `status=${list.status}`);
  const c = await j('POST', '/api/admin/platform/announcements', { token: saToken, body: {
    title: `Regression Announcement ${Date.now()}`, description: 'temp', category: 'General',
    audience: 'All Users', priority: 'Low', publish_date: '2026-07-22', expiry_date: '2026-08-22', status: 'draft'
  }});
  const aid = (c.data?.announcement ?? c.data)?.id;
  check('5.15 create announcement (draft)', c.status === 200 || c.status === 201, `status=${c.status} ${JSON.stringify(c.data).slice(0,150)}`);
  if (aid) {
    const p = await j('POST', `/api/admin/platform/announcements/${aid}/publish`, { token: saToken });
    check('5.16 publish announcement', p.status === 200, `status=${p.status}`);
    const g = await j('GET', `/api/admin/platform/announcements/${aid}`, { token: saToken });
    check('5.17 published announcement readable', g.status === 200, `status=${g.status}`);
    const d = await j('DELETE', `/api/admin/platform/announcements/${aid}`, { token: saToken });
    check('5.18 delete announcement', d.status === 200, `status=${d.status}`);
  } else {
    check('5.16 publish announcement', false, 'no id'); check('5.17 published announcement readable', false, 'no id'); check('5.18 delete announcement', false, 'no id');
  }
}

console.log('=== 6. MISC ====================================================');
{
  const r = await j('GET', '/api/nonexistent');
  check('6.1 unknown endpoint → 404 JSON', r.status === 404, `status=${r.status}`);
}

console.log('\n================================================================');
console.log(`TOTAL: ${passed + failed}   PASSED: ${passed}   FAILED: ${failed}`);
if (failures.length) { console.log('\nFailures:'); failures.forEach(f => console.log('  - ' + f)); process.exitCode = 1; }
};

run().catch(e => { console.error('SUITE CRASHED:', e); process.exit(2); });
