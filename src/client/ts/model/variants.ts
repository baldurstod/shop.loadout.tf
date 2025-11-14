import { VariantJSON } from '../responses/variant';
import { Variant } from './variant';

export class Variants {
	#variants: Variant[] = [];

	add(variant: Variant): void {
		this.#variants.push(variant);
	}

	get count(): number {
		return this.#variants.length;
	}

	[Symbol.iterator](): Iterator<Variant> {
		let index = -1;
		const variants = this.#variants;

		return {
			next: () => ({ value: variants[++index]!, done: !(index in variants) })
		};
	};

	fromJSON(shopVariantsJson: VariantJSON[] = []): void {
		this.#variants = [];

		for (const shopVariantJson of shopVariantsJson) {
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
