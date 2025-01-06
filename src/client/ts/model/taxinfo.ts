import { JSONObject } from '../types';

export class TaxInfo {
	#required: boolean = false;
	#rate: number = 0;
	#shippingTaxable: boolean = false;

	fromJSON(json: JSONObject) {
		this.#required = json.required as boolean;
		this.#rate = Number(json.rate);
		this.#shippingTaxable = json.shipping_taxable as boolean;
	}

	get rate() {
		return this.#rate;
	}

	get shippingTaxable() {
		return this.#shippingTaxable;
	}

	toJSON() {
		return {
			required: this.#required,
			rate: this.#rate,
			shipping_taxable: this.#shippingTaxable
		}
	}
}
