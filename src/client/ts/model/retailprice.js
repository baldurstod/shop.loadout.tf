export class RetailPrice {
	#prices = new Map()
	constructor() {
	}

	setPrice(currency, price) {
		this.#prices.set(currency, price);
	}

	deletePrice(currency) {
		this.#prices.delete(currency);
	}

	getPrice(currency) {
		return this.#prices.get(currency);
	}

	fromJSON(json) {
		this.#prices.clear();
		if (!json) {
			return;
		}

		for (let currency in json) {
			this.setPrice(currency, json[currency]);
		}
	}

	toJSON() {
		const pricesJSON = {};
		for (let [currency, price] of this.#prices) {
			pricesJSON[currency] = price;
		}

		return pricesJSON;
	}
}
