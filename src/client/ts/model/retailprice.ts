import { JSONObject } from 'harmony-types';;

export class RetailPrice {
	#prices = new Map<string, number>()

	setPrice(currency: string, price: number) {
		this.#prices.set(currency, price);
	}

	deletePrice(currency: string) {
		this.#prices.delete(currency);
	}

	getPrice(currency: string) {
		return this.#prices.get(currency);
	}

	fromJSON(json: JSONObject) {
		this.#prices.clear();
		if (!json) {
			return;
		}

		for (let currency in json) {
			this.setPrice(currency, json[currency] as number);
		}
	}

	toJSON() {
		const pricesJSON: JSONObject = {};
		for (let [currency, price] of this.#prices) {
			pricesJSON[currency] = price;
		}

		return pricesJSON;
	}
}
