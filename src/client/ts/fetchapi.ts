import { JSONObject } from "./types"

export async function fetchApi(body: any): Promise<{ requestId: string, response: JSONObject }> {
	const requestId = crypto.randomUUID();
	const response = await fetch('/api', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'X-Request-ID': requestId,
		},
		body: JSON.stringify(body),
	});

	return { requestId: requestId, response: await response.json() };
}
