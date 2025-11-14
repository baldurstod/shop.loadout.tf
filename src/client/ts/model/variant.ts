import { JSONObject } from 'harmony-types';
import { OptionJSON } from '../responses/option';
import { Options } from './options';

export class Variant {
	#id = '';
	#name = '';
	#thumbnailUrl = '';
	#retailPrice = 0;
	#currency = '';
	#options = new Options();

	get id(): string {
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

	fromJSON(shopProductJson: JSONObject = {}): void {
		this.id = shopProductJson.id as string;
		this.name = shopProductJson.name as string;
		this.thumbnailUrl = shopProductJson.thumbnail_url as string;
		this.retailPrice = shopProductJson.retail_price as number;
		this.#currency = shopProductJson.currency as string;
		this.#options.fromJSON(shopProductJson.options as OptionJSON[]);
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
