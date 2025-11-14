import { JSONObject } from 'harmony-types';

export class OrderItem {
	#productId = '';
	#name = '';
	#quantity = 0;
	#retailPrice = 0;
	#thumbnailUrl = '';

	setName(name: string): void {
		this.#name = name;
	}

	getName(): string {
		return this.#name;
	}

	setRetailPrice(retailPrice: number): void {
		this.#retailPrice = retailPrice;
	}

	getRetailPrice(): number {
		return this.#retailPrice;
	}

	setQuantity(quantity: number) {
		this.#quantity = Math.round(quantity);
	}

	getQuantity(): number {
		return this.#quantity;
	}

	setThumbnailUrl(thumbnailUrl: string): void {
		this.#thumbnailUrl = thumbnailUrl;
	}

	getThumbnailUrl(): string {
		return this.#thumbnailUrl;
	}

	getSubtotal(): number {
		return this.#quantity * this.#retailPrice;
	}

	fromJSON(json: JSONObject): void {
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
