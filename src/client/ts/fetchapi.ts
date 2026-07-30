import { ApiResponse } from './responses/response';

export async function fetchApi<T>(action: string, version: number, params: object = {}): Promise<{ requestId: string, response: ApiResponse<T> }> {
	const requestId = crypto.randomUUID();
	const response = await fetch('/api', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json',
			'X-Request-ID': requestId,
		},
		body: JSON.stringify(
			{
				action: action,
				version: version,
				params: params,
			}
		),
	});

	return { requestId: requestId, response: await response.json() as ApiResponse<T> };
}
