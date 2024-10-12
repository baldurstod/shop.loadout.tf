//import { CartItem } from './cartitem.js';
import { DEFAULT_CURRENCY, MAX_PRODUCT_QTY } from '../constants';

export class Cart {
	#items = new Map();
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

	forEach(callbackFn) {
		for (let [_, product] of this.#items) {
			callbackFn(product);
		}
	}

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
		if (this.#items.has(productId)) {
			if (quantity == 0) {
				this.#items.delete(productId);
			} else {
				this.#items.get(productId).setQuantity(quantity);
			}
		}
	}

	clear() {
		this.#items.clear();
	}

	fromJSON(cartJSON) {
		this.#items.clear();
		if (!cartJSON) {
			return;
		}

		const items = cartJSON.items;
		for (let productId in items) {
			const quantity = items[productId];
			this.addProduct(productId, quantity);
		}
	}

	toJSON() {
		const items = {};
		for (let [productId, quantity] of this.#items) {
			items[productId] = quantity;
		}

		return {
			currency: this.currency,
			items: items,
		}
	}
}
