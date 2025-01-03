export class ProductVariant {
	#productId;
	#variantId;
	#name;
	#image;
	constructor() {
	}

	get productId() {
		return this.#productId;
	}

	set productId(productId) {
		this.#productId = productId;
	}

	get variantId() {
		return this.#variantId;
	}

	set variantId(variantId) {
		this.#variantId = variantId;
	}

	get name() {
		return this.#name;
	}

	set name(name) {
		this.#name = name;
	}

	get image() {
		return this.#image;
	}

	set image(image) {
		this.#image = image;
	}

	fromJSON(productJson = {}) {
		this.#productId = productJson.productId;
		this.#variantId = productJson.variantId;
		this.#name = productJson.name;
		this.#image = productJson.image;
	}

	toJSON() {
		return {
			productId: this.productId,
			variantId: this.variantId,
			name: this.name,
			image: this.image,
		}
	}
}
