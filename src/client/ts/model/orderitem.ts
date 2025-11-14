import { OrderItemJSON } from '../responses/order';

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

	setQuantity(quantity: number): void {
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

	fromJSON(json: OrderItemJSON): void {
		this.#productId = json.product_id;
		this.#name = json.name;
		this.setQuantity(json.quantity);
		this.setRetailPrice(Number(json.retail_price));
		this.#thumbnailUrl = json.thumbnail_url;
	}

	toJSON(): OrderItemJSON {
		return {
			product_id: this.#productId,
			name: this.#name,
			quantity: this.#quantity,
			retail_price: String(this.#retailPrice),
			thumbnail_url: this.#thumbnailUrl,
		}
	}
}
