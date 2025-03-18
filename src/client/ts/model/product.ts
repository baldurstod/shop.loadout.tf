import { Files } from './files';
import { Options } from './options';
import { Variants } from './variants';
import { formatPrice } from '../utils';
import { OptionType } from './option';
import { Variant } from './variant';

const productPrices = new Map<string, Map<string, string>>();

export function setRetailPrice(currency: string, productID: string, retailPrice: string) {
	if (!productPrices.has(currency)) {
		productPrices.set(currency, new Map<string, string>());
	}
	productPrices.get(currency)?.set(productID, retailPrice);
}


export function getRetailPrice(currency: string, productID: string): string | undefined {
	return productPrices.get(currency)?.get(productID);
}

export class Product {
	#id: string = '';
	#name: string = '';
	#productName: string = '';
	#thumbnailUrl: string = '';
	#description: string = '';
	#isIgnored = false;
	#dateCreated = Date.now();
	#dateModified = Date.now();
	//#retailPrice: number = 0;
	//#currency: string = '';
	#files = new Files();
	#variantIds = [];
	#externalVariantId: string = '';
	#hasMockupPictures = false;
	#options = new Options();
	#variants = new Variants();
	#status = 'created';

	getId() {
		return this.#id;
	}

	setId(id: string) {
		this.#id = id;
	}

	get externalVariantId() {
		return this.#externalVariantId;
	}

	set externalVariantId(externalVariantId) {
		this.#externalVariantId = externalVariantId;
	}

	get name() {
		return this.#name;
	}

	set name(name) {
		this.#name = name;
	}

	get productName() {
		return this.#productName;
	}

	set productName(productName) {
		this.#productName = productName;
	}

	get thumbnailUrl() {
		return this.#thumbnailUrl;
	}

	set thumbnailUrl(thumbnailUrl) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	get description() {
		return this.#description;
	}

	set description(description) {
		this.#description = description;
	}

	get isIgnored() {
		return this.#isIgnored;
	}

	set isIgnored(isIgnored) {
		this.#isIgnored = isIgnored;
	}

	get status() {
		return this.#status;
	}

	set status(status) {
		this.#status = status;
	}

	get dateCreated() {
		return this.#dateCreated;
	}

	set dateCreated(dateCreated) {
		this.#dateCreated = dateCreated;
	}

	get dateModified() {
		return this.#dateModified;
	}

	set dateModified(dateModified) {
		this.#dateModified = dateModified;
	}

	getRetailPrice(currency: string) {
		return Number(getRetailPrice(currency, this.#id));
	}
	/*
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
	*/

	get options() {
		return this.#options;
	}

	set options(options) {
		this.#options = options;
	}

	get variants() {
		return this.#variants;
	}

	getVariants() {
		return this.#variants;
	}

	get hasMockupPictures() {
		return this.#hasMockupPictures;
	}

	set hasMockupPictures(hasMockupPictures) {
		this.#hasMockupPictures = hasMockupPictures;
	}

	get files() {
		return this.#files;
	}

	set files(files) {
		throw 'remove me';
		this.#files = files;
	}

	get variantIds() {
		return this.#variantIds;
	}

	set variantIds(variantIds) {
		this.#variantIds = variantIds;
	}

	get images() {
		return this.#files.images;
	}

	getThumbnailUrl(fileType: string) {
		return this.#files.getThumbnailUrl(fileType);
	}

	getPriceRange(currency: string) {
		let min = Infinity;
		let max = 0;

		for (let shopVariant of this.#variantIds) {
			const retailPrice = getRetailPrice(currency, shopVariant);
			if (retailPrice === undefined) {
				continue;
			}

			const price = Number(retailPrice);
			min = Math.min(min, price);
			max = Math.max(max, price);
		}

		if (min == Infinity) {
			min = 0;
		}

		return { min: min, max: max, currency: currency };
	}

	addOption(name: string, type: OptionType, value: any) {
		this.#options.addOption(name, type, value);
	}

	addFile(type: string, url: string) {
		this.#files.addFile(type, url);
	}

	addVariant(shopVariant: Variant) {
		this.#variants.add(shopVariant);
	}

	fromJSON(shopProductJson: any = {}) {
		this.#id = shopProductJson.id;
		this.externalVariantId = shopProductJson.external_variant_id;
		this.name = shopProductJson.name;
		this.productName = shopProductJson.product_name;
		this.thumbnailUrl = shopProductJson.thumbnail_url;
		this.#description = shopProductJson.description;
		this.isIgnored = shopProductJson.is_ignored;
		this.#status = shopProductJson.status;
		this.#dateCreated = shopProductJson.date_created;
		this.#dateModified = shopProductJson.date_modified;
		/*
		this.retailPrice = shopProductJson.retail_price;
		this.#currency = shopProductJson.currency;
		*/
		this.#options.fromJSON(shopProductJson.options);
		this.#variants.fromJSON(shopProductJson.variants);
		this.#hasMockupPictures = shopProductJson.has_mockup_pictures;
		this.#files.fromJSON(shopProductJson.files);
		this.#variantIds = shopProductJson.variant_ids;
	}

	toJSON() {
		return {
			_id: this.#id,
			id: this.#id,
			externalVariantId: this.externalVariantId,
			name: this.name,
			productName: this.productName,
			thumbnailUrl: this.thumbnailUrl,
			description: this.#description,
			isIgnored: this.isIgnored,
			status: this.#status,
			dateCreated: this.#dateCreated,
			dateModified: this.#dateModified,
			/*
			retailPrice: this.retailPrice,
			currency: this.#currency,
			*/
			options: this.#options.toJSON(),
			variants: this.#variants.toJSON(),
			hasMockupPictures: this.#hasMockupPictures,
			files: this.#files.toJSON(),
			variantIds: this.#variantIds,
		};
	}

	formatPrice(currency: string): string {
		const retailPrice = getRetailPrice(currency, this.#id);
		if (retailPrice) {
			return formatPrice(Number(retailPrice), currency);
		}
		return "";
	}
}
