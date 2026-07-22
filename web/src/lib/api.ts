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
