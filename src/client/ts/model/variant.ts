import { Options } from './options';

export class Variant {
	#id:string = '';
	#name:string = '';
	#thumbnailUrl:string = '';
	#retailPrice:number = 0;
	#currency:string = '';
	#options = new Options();

	get id() {
		return this.#id;
	}

	set id(id) {
		this.#id = id;
	}

	get name(): string {
		return this.#name;
	}

	set name(name) {
		this.#name = name;
	}

	get thumbnailUrl(): string {
		return this.#thumbnailUrl;
	}

	set thumbnailUrl(thumbnailUrl) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	get retailPrice(): number {
		return this.#retailPrice;
	}

	set retailPrice(retailPrice) {
		this.#retailPrice = Number(retailPrice);
	}

	get currency(): string {
		return this.#currency;
	}

	set currency(currency) {
		this.#currency = currency;
	}

	get options(): Options {
		return this.#options;
	}

	set options(options) {
		this.#options = options;
	}

	fromJSON(shopProductJson: any = {}) {
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
