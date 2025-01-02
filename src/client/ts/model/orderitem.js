import {CartProductVariant} from './cartproductvariant.js';

export class OrderItem {
	#files;
	#quantity;
	#retailPrice;
	#thumbnailUrl;
	constructor() {
		this.id;
		this.externalId;
		this.variantId;
		this.syncVariantId;
		this.externalVariantId;
		this.warehouseProductVariantId;
		this.#quantity;
		this.price;
		this.#retailPrice;
		this.name;
		this.product = new CartProductVariant();
		this.#files = new Set();
		this.options;
		this.sku;
	}

	set files(files) {
		for (let file of files) {
			this.#files.add(file);
		}
	}

	getFileUrl(fileType) {
		for (let file of this.#files) {
			if (file.type == fileType) {
				return file;
			}
		}
	}

	set retailPrice(retailPrice) {
		this.#retailPrice = Number.parseFloat(retailPrice);
	}

	get retailPrice() {
		return this.#retailPrice;
	}

	set quantity(quantity) {
		this.#quantity = Math.round(quantity);
	}

	get quantity() {
		return this.#quantity;
	}

	set thumbnailUrl(thumbnailUrl) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	get thumbnailUrl() {
		return this.#thumbnailUrl;
	}

	get subtotal() {
		return this.#quantity * this.#retailPrice;
	}

	fromJSON(json) {
		this.externalVariantId = json.externalVariantId;
		this.name = json.name;
		this.quantity = json.quantity;
		this.retailPrice = json.retail_price;
		this.#thumbnailUrl = json.thumbnail_url;

		this.#files = [];
		if (json.files) {
			/*json.files.forEach(
				(jsonFile) => {
					let file = new File();
					file.fromJSON(jsonFile);
					this.#files.push(file);
				}
			);*/
		}
	}

	toJSON() {
		return {
			externalVariantId: this.externalVariantId,
			name: this.name,
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
