import { RetailPrice } from '../retailprice.js'
import { mapToObject, objectToMap } from '../../utils/map'

export class PrintfulProduct {
	#id;
	#title;
	#description;
	#brand;
	#model;
	#image;
	#retailPrice = new RetailPrice();
	#files = new Set();
	#images = new Set();
	//#options = new ProductOptions();
	#attributes = new Map();
	#variants = new Set();
	#customData;

	constructor() {
	}

	get id() {
		return this.#id;
	}

	set id(id) {
		this.#id = id;
	}

	get title() {
		return this.#title;
	}

	set title(title) {
		this.#title = title;
	}

	get brand() {
		return this.#brand;
	}

	set brand(brand) {
		this.#brand = brand;
	}

	get model() {
		return this.#model;
	}

	set model(model) {
		this.#model = model;
	}

	get image() {
		return this.#image;
	}

	set image(image) {
		this.#image = image;
	}

	get files() {
		return this.#files;
	}

	set files(files) {
		this.#files = files;
	}

	/*get options() {
		return this.#options;
	}

	set options(options) {
		this.#options = options;
	}*/

	get description() {
		return this.#description;
	}

	set description(description) {
		this.#description = description;
	}

	get retailPrice() {
		return this.#retailPrice;
	}

	setRetailPrice(currency, price) {
		this.#retailPrice.setPrice(currency, price);
	}

	getRetailPrice(currency) {
		return this.#retailPrice.getPrice(currency);
	}

	setAttribute(attributeName, attributeValue) {
		this.#attributes.set(attributeName, attributeValue);
	}

	deleteAttribute(attributeName) {
		this.#attributes.delete(attributeName);
	}

	addFile(file) {
		this.#files.add(file);
	}

	removeFile(file) {
		this.#files.delete(file);
	}

	addImage(image) {
		this.#images.add(image);
	}

	removeImage(image) {
		this.#images.delete(image);
	}

	addVariant(id) {
		this.#variants.add(id);
	}

	removeVariant(id) {
		this.#variants.delete(id);
	}

	get customData() {
		return this.#customData;
	}

	set customData(customData) {
		this.#customData = customData;
	}

	fromJSON(productJson: any = {}) {
		this.id = productJson.id;
		this.title = productJson.title;
		this.brand = productJson.brand;
		this.model = productJson.model;
		this.image = productJson.image;
		//this.options = productJson.options;
		this.description = productJson.description;
		this.#retailPrice.fromJSON(productJson.retailPrice);

		this.#variants.clear();
		if (productJson.variants) {
			productJson.variants.map(item => this.addVariant(item));
		}

		this.#images.clear();
		if (productJson.images) {
			productJson.images.map(image => this.addImage(image));
		}


		objectToMap(productJson.attributes, this.#attributes);
		return;
		/*

		console.log(this.#files);

		this.#files.clear();
		for (let f of productJson.files) {
			const file = new File();
			file.fromJSON(f);
			this.#files.add(file);
		}

		this.customData = productJson.customData;
		*/
	}

	toJSON() {
		return {
			id: this.id,
			title: this.title,
			brand: this.brand,
			model: this.model,
			image: this.image,
			files: this.files,
			images: Array.from(this.#images),
			//options: this.options,
			description: this.description,
			retailPrice: this.retailPrice.toJSON(),
			attributes: mapToObject(this.#attributes),
			variants: Array.from(this.#variants),
			customData: this.customData,
		}
	}
}
