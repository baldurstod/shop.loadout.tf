import { Variant } from './variant';

export class Variants {
	#variants: Array<Variant> = [];

	add(variant: Variant) {
		this.#variants.push(variant);
	}

	get count() {
		return this.#variants.length;
	}

	[Symbol.iterator]() {
		let index = -1;
		let variants = this.#variants;

		return {
			next: () => ({ value: variants[++index], done: !(index in variants) })
		};
	};

	fromJSON(shopVariantsJson = []) {
		this.#variants = [];

		for (let shopVariantJson of shopVariantsJson) {
			const shopVariant = new Variant();
			shopVariant.fromJSON(shopVariantJson);
			this.#variants.push(shopVariant);
		}
	}

	toJSON() {
		const variants = [];
		for (const variant of this.#variants) {
			variants.push(variant.toJSON());
		}
		return variants;
	}
}
