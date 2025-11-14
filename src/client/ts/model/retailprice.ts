import { JSONObject } from 'harmony-types';

export class RetailPrice {
	#prices = new Map<string, number>()

	setPrice(currency: string, price: number): void {
		this.#prices.set(currency, price);
	}

	deletePrice(currency: string): void {
		this.#prices.delete(currency);
	}

	getPrice(currency: string): number | null {
		return this.#prices.get(currency) ?? null;
	}

	fromJSON(json: JSONObject): void {
		this.#prices.clear();
		if (!json) {
			return;
		}

		for (const currency in json) {
			this.setPrice(currency, json[currency] as number);
		}
	}

	toJSON() {
		const pricesJSON: JSONObject = {};
		for (const [currency, price] of this.#prices) {
			pricesJSON[currency] = price;
		}

		return pricesJSON;
	}
}
