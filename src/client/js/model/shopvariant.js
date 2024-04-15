import { ShopOptions } from './shopoptions.js';

export class ShopVariant {
	#id;
	#name;
	#thumbnailUrl;
	#retailPrice;
	#currency;
	#options = new ShopOptions();

	constructor() {
	}

	get id() {
		return this.#id;
	}

	set id(id) {
		this.#id = id;
	}

	get name() {
		return this.#name;
	}

	set name(name) {
		this.#name = name;
	}

	get thumbnailUrl() {
		return this.#thumbnailUrl;
	}

	set thumbnailUrl(thumbnailUrl) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	get retailPrice() {
		return this.#retailPrice;
	}

	set retailPrice(retailPrice) {
		this.#retailPrice = Number(retailPrice);
	}

	get currency() {
		return this.#currency;
	}

	set currency(currency) {
		this.#currency = currency;
	}

	get options() {
		return this.#options;
	}

	set options(options) {
		this.#options = options;
	}

	fromJSON(shopProductJson = {}) {
		this.id = shopProductJson.id;
		this.name = shopProductJson.name;
		this.thumbnailUrl = shopProductJson.thumbnail_url;
		this.retailPrice = shopProductJson.retail_price;
		this.#currency = shopProductJson.currency;
		this.#options.fromJSON(shopProductJson.options);
	}

	toJSON() {
		return {
			id: this.id,
			name: this.name,
			thumbnail_url: this.thumbnailUrl,
			retail_price: this.retailPrice,
			currency: this.#currency,
			options: this.#options.toJSON(),
		};
	}
}
