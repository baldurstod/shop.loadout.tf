import { File } from '../file.js';
import { ProductVariant } from '../productvariant.js';

import { formatPrice } from '../../utils.js';

export class SyncVariant {
	#syncProduct;
	#id;
	#externalId;
	#syncProductId;
	#name;
	#synced;
	#variantId;
	#retailPrice;
	#currency;
	#isIgnored;
	#sku;
	#product;
	#files;
	#options;
	#warehouseProductVariantId;
	#lastUpdate = 0;

	constructor(syncProduct) {
		this.#syncProduct = syncProduct;
	}

	get syncProduct() {
		return this.#syncProduct;
	}

	set syncProduct(syncProduct) {
		this.#syncProduct = syncProduct;
	}

	get id() {
		return this.#id;
	}

	set id(id) {
		this.#id = id;
	}

	get externalId() {
		return this.#externalId;
	}

	set externalId(externalId) {
		this.#externalId = externalId;
	}

	get syncProductId() {
		return this.#syncProductId;
	}

	set syncProductId(syncProductId) {
		this.#syncProductId = syncProductId;
	}

	get name() {
		return this.#name;
	}

	set name(name) {
		this.#name = name;
	}

	get synced() {
		return this.#synced;
	}

	set synced(synced) {
		this.#synced = synced;
	}

	get variantId() {
		return this.#variantId;
	}

	set variantId(variantId) {
		this.#variantId = variantId;
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

	get isIgnored() {
		return this.#isIgnored;
	}

	set isIgnored(isIgnored) {
		this.#isIgnored = isIgnored;
	}

	get sku() {
		return this.#sku;
	}

	set sku(sku) {
		this.#sku = sku;
	}

	get product() {
		return this.#product;
	}

	set product(product) {
		this.#product = product;
	}

	get files() {
		return this.#files;
	}

	set files(files) {
		this.#files = files;
	}

	get options() {
		return this.#options;
	}

	set options(options) {
		this.#options = options;
	}

	get warehouseProductVariantId() {
		return this.#warehouseProductVariantId;
	}

	set warehouseProductVariantId(warehouseProductVariantId) {
		this.#warehouseProductVariantId = warehouseProductVariantId;
	}

	get lastUpdate() {
		return this.#lastUpdate;
	}

	set lastUpdate(lastUpdate) {
		this.#lastUpdate = lastUpdate;
	}

	fromJSON(syncVariantJson: any = {}) {
		this.id = syncVariantJson.id;
		this.externalId = syncVariantJson.externalId;
		this.syncProductId = syncVariantJson.syncProductId;
		this.name = syncVariantJson.name;
		this.synced = syncVariantJson.synced;
		this.variantId = syncVariantJson.variantId;
		this.retailPrice = syncVariantJson.retailPrice;
		this.currency = syncVariantJson.currency;
		this.isIgnored = syncVariantJson.isIgnored;
		this.sku = syncVariantJson.sku;
		this.product = new ProductVariant();
		this.product.fromJSON(syncVariantJson.product);
		this.#lastUpdate = syncVariantJson.lastUpdate ?? 0;

		this.#files = [];

		let filesJson = syncVariantJson.files;
		for (let fileJson of filesJson) {
			let file = new File();
			file.fromJSON(fileJson);
			this.#files.push(file);
		}
		this.options = syncVariantJson.options;
		this.warehouseProductVariantId = syncVariantJson.warehouseProductVariantId;
	}

	toJSON() {
		const files = [];
		for (let file of this.files) {
			files.push(file.toJSON());
		}

		return {
			id: this.id,
			externalId: this.externalId,
			syncProductId: this.syncProductId,
			name: this.name,
			synced: this.synced,
			variantId: this.variantId,
			retailPrice: this.retailPrice,
			currency: this.currency,
			isIgnored: this.isIgnored,
			sku: this.sku,
			product: this.product.toJSON(),
			files: files,
			options: this.options,
			warehouseProductVariantId: this.warehouseProductVariantId,
			lastUpdate: this.#lastUpdate,
		}
	}

	get images() {
		let images = []
		for (let file of this.files) {
			images.push(file.previewUrl ?? file.url);
		}
		return images;
	}

	getThumbnailUrl(fileType) {
		for (let file of this.files) {
			if (file.type == fileType) {
				return file.thumbnailUrl;
			}
		}
	}

	get externalProductId() {
		return this.#syncProduct.externalId;
	}

	formatPrice() {
		return formatPrice(this.retailPrice, this.currency);
	}
}
