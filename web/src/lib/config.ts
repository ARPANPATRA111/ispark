import { env } from '$env/dynamic/public';

// Values configured through a hosting dashboard can pick up stray whitespace or
// a byte-order mark. A BOM in front of the URL makes it fail `new URL()`, so
// fetch treats it as a relative path and every API call quietly hits the web
// origin instead, returning HTML that then breaks JSON parsing. Strip those
// characters (and any trailing slash) so a copy-paste artefact cannot take the
// whole frontend down.
//
// U+FEFF (byte-order mark) and U+200B (zero-width space) are built from their
// code points rather than written literally: they are invisible in source too,
// which is exactly how this class of bug hides.
const INVISIBLE_CHARS = String.fromCharCode(0xfeff, 0x200b);
const INVISIBLE_PREFIX = new RegExp(`^[${INVISIBLE_CHARS}\\s]+`);
const TRAILING_JUNK = /[\s/]+$/;

function sanitizeBaseUrl(value: string | undefined, fallback: string): string {
	const cleaned = (value ?? '').replace(INVISIBLE_PREFIX, '').replace(TRAILING_JUNK, '');
	if (!cleaned) return fallback;

	try {
		new URL(cleaned);
		return cleaned;
	} catch {
		console.error(
			`Invalid PUBLIC_API_BASE_URL (${JSON.stringify(value)}); falling back to`,
			fallback
		);
		return fallback;
	}
}

export const API_BASE_URL = sanitizeBaseUrl(env.PUBLIC_API_BASE_URL, 'http://localhost:8080');
