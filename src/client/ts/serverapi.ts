export class ServerAPI {
	static async #fetchAPI(action, version, params: any = {}) {
		let fetchOptions = {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify(
				{
					action: action,
					version: version,
					params: params,
				}
			),
		};

		const response = await fetch('./api', fetchOptions);
		const json = await response.json();
		if (json.success) {
			return json.result;
		}

		throw json.error;
	}

	static async getCountries() {
		const result = await this.#fetchAPI('get-countries', 1);
		return result?.countries;
	}

	static async getCurrency() {
		return this.#fetchAPI('get-currency', 1);
	}
}
/*
export class ServerAPI {
	static async #fetchAPI(action, version, params) {
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
		return { requestId: requestId, response: await response.json() };
	}

	static async getCountries() {
		const { requestId, response } = this.#fetchAPI('get-countries', 1);
		if (response.success) {
			return response.result?.countries;
		}
	}

	static async getCurrency() {
		const { requestId, response } = this.#fetchAPI('get-currency', 1);
		if (response.success) {
			return json.result;
		}
	}
}
*/
