import { JSONObject } from '../types';

export class OrderItem {
	#productId: string = '';
	#name: string = '';
	#quantity: number = 0;
	#retailPrice: number = 0;
	#thumbnailUrl: string = '';

	setName(name: string) {
		this.#name = name;
	}

	getName() {
		return this.#name;
	}

	setRetailPrice(retailPrice: number) {
		this.#retailPrice = retailPrice;
	}

	getRetailPrice() {
		return this.#retailPrice;
	}

	setQuantity(quantity: number) {
		this.#quantity = Math.round(quantity);
	}

	getQuantity() {
		return this.#quantity;
	}

	setThumbnailUrl(thumbnailUrl: string) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	getThumbnailUrl() {
		return this.#thumbnailUrl;
	}

	getSubtotal() {
		return this.#quantity * this.#retailPrice;
	}

	fromJSON(json: JSONObject) {
		this.#name = json.name as string;
		this.setQuantity(json.quantity as number);
		this.setRetailPrice(json.retail_price as number);
		this.#thumbnailUrl = json.thumbnail_url as string;
	}

	toJSON(): JSONObject {
		return {
			name: this.#name,
			quantity: this.#quantity,
			retailPrice: this.#retailPrice,
			thumbnailUrl: this.#thumbnailUrl,
		}
	}
}

/*id	{…}
external_id	{…}
variant_id	{…}
sync_variant_id	{…}
external_variant_id	{…}
warehouse_product_variant_id	{…}
quantity	{…}
price	{…}
retail_price	{…}
name	{…}
product	{…}
files	{…}
options	{…}
sku	{…}*/
