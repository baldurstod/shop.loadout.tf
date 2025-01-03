import { JSONObject } from '../types';

export class ShippingInfo {
	id: string = '';
	name: string = '';
	rate: number = 0;
	currency: string = '';
	minDeliveryDays: number = 0;
	maxDeliveryDays: number = 0;

	fromJSON(json: JSONObject) {
		this.id = json.id as string;
		this.name = json.name as string;
		this.rate = Number(json.rate);
		this.currency = json.currency as string;
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
		}
	}
}
