import { JSONObject } from 'harmony-types';
import { DEFAULT_CURRENCY } from '../constants';
import { CartJSON } from '../responses/cart';
;

export class Cart {
	#items = new Map<string, number>();
	#currency = DEFAULT_CURRENCY;

	set currency(currency) {
		this.#currency = currency;
		this.#items.clear();
	}

	get currency(): string {
		return this.#currency;
	}

	get items(): Map<string, number> {
		return this.#items;
	}

	/*
	forEach(callbackFn: (product: Product) => void) {
		for (let [_, product] of this.#items) {
			callbackFn(product);
		}
	}
	*/

	get totalQuantity(): number {
		let quantity = 0;
		for (const [, qty] of this.#items) {
			quantity += qty;
		}
		return quantity;
	}

	addProduct(productId: string, quantity: number): void {
		this.#items.set(productId, quantity);
	}

	setQuantity(productId: string, quantity: number): void {
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

	clear(): void {
		this.#items.clear();
	}

	fromJSON(cart: CartJSON): void {
		this.#items.clear();
		if (!cart) {
			return;
		}

		const items = cart.items as JSONObject;
		for (const productId in items) {
			const quantity = items[productId];
			this.addProduct(productId, quantity as number);
		}
	}

	toJSON() {
		const items: JSONObject = {};
		for (const [productId, quantity] of this.#items) {
			items[productId] = quantity;
		}

		return {
			currency: this.currency,
			items: items,
		}
	}
}
