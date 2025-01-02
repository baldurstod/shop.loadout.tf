export class TaxInfo {
	#required;
	#rate;
	#shippingTaxable;
	constructor() {
	}

	fromJSON(json) {
		this.#required = json.required;
		this.#rate = Number(json.rate);
		this.#shippingTaxable = json.shippingTaxable;
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
			shippingTaxable: this.#shippingTaxable
		}
	}
}
