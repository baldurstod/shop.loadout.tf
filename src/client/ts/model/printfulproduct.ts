export class PrintfulProduct {
	#id;
	#type;
	#typeName;
	#title;
	#brand;
	#model;
	#image;
	#variantCount;
	#currency;
	#files;
	#options;
	#isDiscontinued;
	#avgFulfillmentTime;
	#description;

	get id() {
		return this.#id;
	}

	set id(id) {
		this.#id = id;
	}

	get type() {
		return this.#type;
	}

	set type(type) {
		this.#type = type;
	}

	get typeName() {
		return this.#typeName;
	}

	set typeName(typeName) {
		this.#typeName = typeName;
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

	get variantCount() {
		return this.#variantCount;
	}

	set variantCount(variantCount) {
		this.#variantCount = variantCount;
	}

	get currency() {
		return this.#currency;
	}

	set currency(currency) {
		this.#currency = currency;
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

	get isDiscontinued() {
		return this.#isDiscontinued;
	}

	set isDiscontinued(isDiscontinued) {
		this.#isDiscontinued = isDiscontinued;
	}

	get avgFulfillmentTime() {
		return this.#avgFulfillmentTime;
	}

	set avgFulfillmentTime(avgFulfillmentTime) {
		this.#avgFulfillmentTime = avgFulfillmentTime;
	}

	get description() {
		return this.#description;
	}

	set description(description) {
		this.#description = description;
	}

	fromJSON(productJson = {}) {
		this.id = productJson.id;
		this.type = productJson.type;
		this.typeName = productJson.typeName;
		this.title = productJson.title;
		this.brand = productJson.brand;
		this.model = productJson.model;
		this.image = productJson.image;
		this.variantCount = productJson.variantCount;
		this.currency = productJson.currency;
		this.files = productJson.files;
		this.options = productJson.options;
		this.isDiscontinued = productJson.isDiscontinued;
		this.avgFulfillmentTime = productJson.avgFulfillmentTime;
		this.description = productJson.description;
	}

	toJSON() {
		return {
			id: this.id,
			type: this.type,
			typeName: this.typeName,
			title: this.title,
			brand: this.brand,
			model: this.model,
			image: this.image,
			variantCount: this.variantCount,
			currency: this.currency,
			files: this.files,
			options: this.options,
			isDiscontinued: this.isDiscontinued,
			avgFulfillmentTime: this.avgFulfillmentTime,
			description: this.description,
		}
	}
}
