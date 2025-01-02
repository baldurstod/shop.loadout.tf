import {MAX_PRODUCT_QTY} from '../../js/constants.js';

export class CartItem {
	#productId;
	#quantity;
	#price;
	constructor(productId, quantity, price) {
		this.#productId = productId;
		this.#quantity = quantity;
		//this.#price = price;
	}

	/*get product() {
		return this.#product;
	}*/

	get productId() {
		return this.#productId;
	}

	get quantity() {
		return this.#quantity;
	}

	/*get price() {
		return this.#price;
	}*/

	get name() {
		throw 'todo get name() {';
		//return this.#name;
	}

	get shopUrl() {
		throw 'todo get shopUrl() {';
		//return this.#shopUrl;
	}

	get thumbnailUrl() {
		throw 'todo get thumbnailUrl() {';
		//return this.#thumbnailUrl;
	}

	/*get subtotal()  {
		return this.#quantity * this.#price;
	}*/

	addQuantity(quantity) {
		this.setQuantity(Math.min(this.#quantity + quantity, MAX_PRODUCT_QTY));
	}

	setQuantity(quantity) {
		this.#quantity = Math.min(quantity, MAX_PRODUCT_QTY);
	}

	toJSON() {
		return { productId: this.#productId, quantity: this.#quantity };
	}

	fromJSON(json) {
		throw 'todo';
		/*this.#product = json.productId;
		this.#quantity = json.quantity;
		this.#price = json.price;
		this.#name = json.name;
		this.#variant = json.variant;
		this.#thumbnailUrl = json.thumbnailUrl;
		this.#shopUrl = json.shopUrl;
		return this;*/
	}
}
