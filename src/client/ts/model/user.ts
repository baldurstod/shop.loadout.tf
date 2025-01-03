import { Address } from './address.js';
import { Cart } from '../../js/model/cart.js';

export class User {
	#id;
	#currency;
	#dateCreated;
	#dateUpdated;
	#lastUsed;
	#shippingAddress;
	#billingAddress;
	#sameBillingAddress;
	#orders;
	#cart;
	#status = 'created';
	constructor(id) {
		this.#id = id;
		this.#currency = 'USD';
		this.#dateCreated = Date.now();
		this.#dateUpdated = Date.now();
		this.#lastUsed = Date.now();
		this.#shippingAddress = new Address();
		this.#billingAddress = new Address();
		this.#sameBillingAddress = true;
		this.#orders = [];
		this.#cart = new Cart();
	}

	get id() {
		return this.#id;
	}

	get shippingAddress() {
		return this.#shippingAddress;
	}

	get billingAddress() {
		return this.#sameBillingAddress ? this.#shippingAddress : this.#billingAddress;
	}

	get sameBillingAddress() {
		return this.#sameBillingAddress;
	}

	set sameBillingAddress(sameBillingAddress) {
		this.#sameBillingAddress = sameBillingAddress;
	}

	set currency(currency) {
		this.#currency = currency;
	}

	get currency() {
		return this.#currency;
	}

	get cart() {
		return this.#cart;
	}

	get dateCreated() {
		return this.#dateCreated;
	}

	set dateCreated(dateCreated) {
		this.#dateCreated = dateCreated;
	}

	get dateUpdated() {
		return this.#dateUpdated;
	}

	get lastUsed() {
		return this.#lastUsed;
	}

	set lastUsed(lastUsed) {
		this.#lastUsed = lastUsed;
	}

	get status() {
		return this.#status;
	}

	set status(status) {
		this.#status = status;
	}

	get orders() {
		return this.#orders;
	}

	addOrder(orderId) {
		this.#orders.push(orderId);
	}

	fromJSON(json) {
		this.#id = json.id;
		this.#currency = json.currency;
		this.#dateCreated = json.dateCreated;
		this.#dateUpdated = json.dateUpdated;
		this.#lastUsed = json.lastUsed;
		this.#status = json.status;
		this.#shippingAddress.fromJSON(json.shippingAddress);
		this.#billingAddress.fromJSON(json.billingAddress);
		this.#sameBillingAddress = json.sameBillingAddress;
		this.#orders = json.orders;
		this.#cart.fromJSON(json.cart);
	}

	toJSON() {
		return {
			id: this.id,
			currency: this.currency,
			dateCreated: this.dateCreated,
			dateUpdated: this.dateUpdated,
			lastUsed: this.lastUsed,
			status: this.status,
			shippingAddress: this.shippingAddress.toJSON(),
			billingAddress: this.billingAddress.toJSON(),
			sameBillingAddress: this.sameBillingAddress,
			orders: this.orders,
			cart: this.#cart.toJSON(),
		}
	}
}
