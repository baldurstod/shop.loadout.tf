import { TaxInfoJSON } from '../responses/order';

export class TaxInfo {
	#required = false;
	#rate = 0;
	#shippingTaxable = false;

	fromJSON(json: TaxInfoJSON): void {
		this.#required = json.required;
		this.#rate = Number(json.rate);
		this.#shippingTaxable = json.shipping_taxable;
	}

	get rate(): number {
		return this.#rate;
	}

	get shippingTaxable(): boolean {
		return this.#shippingTaxable;
	}

	toJSON(): TaxInfoJSON {
		return {
			required: this.#required,
			rate: this.#rate,
			shipping_taxable: this.#shippingTaxable
		}
	}
}
