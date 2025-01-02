export class ShippingInfo {
	constructor(json) {
		this.id;
		this.name;
		this.rate;
		this.currency;
		this.minDeliveryDays;
		this.maxDeliveryDays;
		if (json) {
			this.fromJSON(json);
		}
	}

	fromJSON(json) {
		this.id = json.id;
		this.name = json.name;
		this.rate = Number(json.rate);
		this.currency = json.currency;
		this.minDeliveryDays = Number(json.minDeliveryDays);
		this.maxDeliveryDays = Number(json.maxDeliveryDays);
	}

	toJSON() {
		return {
			id: this.id,
			name: this.name,
			rate: this.rate,
			currency: this.currency,
			minDeliveryDays: this.minDeliveryDays,
			maxDeliveryDays: this.maxDeliveryDays,
			taxInfo: this.taxInfo,
			shippingMethod: this.shippingMethod,
			printfulOrder: this.printfulOrder
		}
	}
}
