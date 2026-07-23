/**
 * Reads a JSON body from an API response, turning the failure modes that are
 * normal for a hosted deployment into messages a user can act on.
 *
 * Calling `response.json()` directly surfaces "Unexpected token '<'" whenever
 * the body is HTML rather than JSON — which happens while a sleeping free-tier
 * instance is waking up, when a gateway returns an error page, or when the API
 * base URL is misconfigured and the request lands on the web origin.
 */
export async function readJson(response: Response): Promise<Record<string, unknown>> {
	const contentType = response.headers.get('content-type') ?? '';

	if (contentType.includes('application/json')) {
		try {
			return (await response.json()) as Record<string, unknown>;
		} catch {
			throw new Error('The server sent a malformed response. Please try again.');
		}
	}

	// Not JSON: report what actually went wrong rather than a parser error.
	if (response.status === 502 || response.status === 503 || response.status === 504) {
		throw new Error('The server is starting up. Please wait a moment and try again.');
	}
	if (!response.ok) {
		throw new Error(`The server returned an unexpected error (${response.status}).`);
	}
	throw new Error('The server sent an unexpected response. Please try again.');
}

/**
 * Re-runs `refresh` when the user returns to the tab (or the window regains
 * focus), so data feels live without polling on a timer.
 *
 * Two guards keep it cheap: it only fires when the document is actually
 * visible, and it ignores triggers that arrive within `minIntervalMs` of the
 * last run. Switching tabs rapidly therefore costs at most one request per
 * interval rather than one per switch.
 *
 * Returns a cleanup function — call it from onMount's teardown.
 */
export function refreshOnFocus(refresh: () => void, minIntervalMs = 30_000): () => void {
	let lastRun = Date.now();

	const maybeRefresh = () => {
		if (document.visibilityState !== 'visible') return;
		const now = Date.now();
		if (now - lastRun < minIntervalMs) return;
		lastRun = now;
		refresh();
	};

	document.addEventListener('visibilitychange', maybeRefresh);
	window.addEventListener('focus', maybeRefresh);

	return () => {
		document.removeEventListener('visibilitychange', maybeRefresh);
		window.removeEventListener('focus', maybeRefresh);
	};
}

/**
 * Fetches a protected file with a bearer token and saves it in the browser.
 * The file endpoints require an Authorization header, so a plain anchor link
 * cannot be used — the bytes are fetched, turned into an object URL, and a
 * temporary link is clicked to trigger the download.
 */
export async function downloadAuthedFile(
	url: string,
	token: string,
	fallbackName = 'certificate'
): Promise<void> {
	const response = await fetch(url, { headers: { Authorization: `Bearer ${token}` } });
	if (!response.ok) {
		// The body is JSON on error; surface its message.
		const data = await readJson(response).catch(() => ({}));
		throw new Error(String((data as { error?: string }).error || 'File is not available.'));
	}

	// Prefer the server-provided filename from Content-Disposition.
	let filename = fallbackName;
	const disposition = response.headers.get('content-disposition') ?? '';
	const match = /filename="?([^"]+)"?/.exec(disposition);
	if (match) filename = match[1];

	const blob = await response.blob();
	const objectUrl = URL.createObjectURL(blob);
	const link = document.createElement('a');
	link.href = objectUrl;
	link.download = filename;
	document.body.appendChild(link);
	link.click();
	link.remove();
	URL.revokeObjectURL(objectUrl);
}
