import { ProductJSON } from '../responses/product';
import { formatPrice } from '../utils';
import { Files } from './files';
import { OptionType } from './option';
import { Options } from './options';
import { Variant } from './variant';
import { Variants } from './variants';

const productPrices = new Map<string, Map<string, string>>();

export function setRetailPrice(currency: string, productID: string, retailPrice: string): void {
	if (!productPrices.has(currency)) {
		productPrices.set(currency, new Map<string, string>());
	}
	productPrices.get(currency)?.set(productID, retailPrice);
}


export function getRetailPrice(currency: string, productID: string): string | undefined {
	return productPrices.get(currency)?.get(productID);
}

export class Product {
	#id = '';
	#name = '';
	#productName = '';
	#thumbnailUrl = '';
	#description = '';
	#isIgnored = false;
	#dateCreated = Date.now();
	#dateModified = Date.now();
	//#retailPrice: number = 0;
	//#currency: string = '';
	#files = new Files();
	#variantIds: string[] = [];
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

	get name(): string {
		return this.#name;
	}

	set name(name) {
		this.#name = name;
	}

	get productName(): string {
		return this.#productName;
	}

	set productName(productName) {
		this.#productName = productName;
	}

	get thumbnailUrl(): string {
		return this.#thumbnailUrl;
	}

	set thumbnailUrl(thumbnailUrl) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	get description(): string {
		return this.#description;
	}

	set description(description) {
		this.#description = description;
	}

	get isIgnored(): boolean {
		return this.#isIgnored;
	}

	set isIgnored(isIgnored) {
		this.#isIgnored = isIgnored;
	}

	get status(): string {
		return this.#status;
	}

	set status(status) {
		this.#status = status;
	}

	get dateCreated(): number {
		return this.#dateCreated;
	}

	set dateCreated(dateCreated) {
		this.#dateCreated = dateCreated;
	}

	get dateModified(): number {
		return this.#dateModified;
	}

	set dateModified(dateModified) {
		this.#dateModified = dateModified;
	}

	getRetailPrice(currency: string): number {
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

	get options(): Options {
		return this.#options;
	}

	set options(options) {
		this.#options = options;
	}

	get variants(): Variants {
		return this.#variants;
	}

	getVariants(): Variants {
		return this.#variants;
	}

	get hasMockupPictures(): boolean {
		return this.#hasMockupPictures;
	}

	set hasMockupPictures(hasMockupPictures) {
		this.#hasMockupPictures = hasMockupPictures;
	}

	get files(): Files {
		return this.#files;
	}

	set files(files) {
		throw new Error('remove me');
		this.#files = files;
	}

	get variantIds(): string[] {
		return this.#variantIds;
	}

	set variantIds(variantIds) {
		this.#variantIds = variantIds;
	}

	get images(): string[] {
		return this.#files.images;
	}

	getThumbnailUrl(fileType: string): string | null {
		return this.#files.getThumbnailUrl(fileType) ?? null;
	}

	getPriceRange(currency: string): { min: number, max: number, currency: string } {
		let min = Infinity;
		let max = 0;

		for (const shopVariant of this.#variantIds) {
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

	fromJSON(shopProductJson: ProductJSON) {
		this.#id = shopProductJson.id;
		this.externalVariantId = shopProductJson.external_variant_id;
		this.name = shopProductJson.name;
		this.productName = shopProductJson.product_name;
		this.thumbnailUrl = shopProductJson.thumbnail_url;
		this.#description = shopProductJson.description;
		this.isIgnored = shopProductJson.is_ignored;
		this.#status = shopProductJson.status;
		this.#dateCreated = shopProductJson.date_created;
		this.#dateModified = shopProductJson.date_updated;
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
