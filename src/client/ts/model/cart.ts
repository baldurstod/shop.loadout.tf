import { DEFAULT_CURRENCY, MAX_PRODUCT_QTY } from '../constants';
import { JSONObject } from '../types';

export class Cart {
	#items = new Map<string, number>();
	#currency = DEFAULT_CURRENCY;

	set currency(currency) {
		this.#currency = currency;
		this.#items.clear();
	}

	get currency() {
		return this.#currency;
	}

	get items() {
		return this.#items;
	}

	/*
	forEach(callbackFn: (product: Product) => void) {
		for (let [_, product] of this.#items) {
			callbackFn(product);
		}
	}
	*/

	get totalQuantity() {
		let quantity = 0;
		for (let [_, qty] of this.#items) {
			quantity += qty;
		}
		return quantity;
	}

	addProduct(productId: string, quantity: number) {
		this.#items.set(productId, quantity);
	}

	setQuantity(productId: string, quantity: number) {
		quantity = Math.floor(quantity);
		if (isNaN(quantity)) {
			return;
		}
		if (quantity < 0) {
			return;
		}
		if (quantity == 0) {
			this.#items.delete(productId);
		} else {
			this.#items.set(productId, quantity);
		}
	}

	clear() {
		this.#items.clear();
	}

	fromJSON(cartJSON: JSONObject) {
		this.#items.clear();
		if (!cartJSON) {
			return;
		}

		const items = cartJSON.items as JSONObject;
		for (let productId in items) {
			const quantity = items[productId];
			this.addProduct(productId, quantity as number);
		}
	}

	toJSON() {
		const items: JSONObject = {};
		for (let [productId, quantity] of this.#items) {
			items[productId] = quantity;
		}

		return {
			currency: this.currency,
			items: items,
		}
	}
}
