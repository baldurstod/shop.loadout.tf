import { JSONArray, JSONObject } from 'harmony-types';

export class ShippingInfo {
	shipping = '';
	shippingMethodName = '';
	rate = ''
	currency = '';
	minDeliveryDays = 0;
	maxDeliveryDays = 0;
	minDeliveryDate = '';
	maxDeliveryDate = '';
	shipments: Shipment[] = [];

	fromJSON(json: JSONObject): void {
		this.shipping = json.shipping as string;
		this.shippingMethodName = json.shipping_method_name as string;
		this.rate = json.rate as string;
		this.currency = json.currency as string;
		this.minDeliveryDays = Number(json.min_delivery_days);
		this.maxDeliveryDays = Number(json.max_delivery_days);
		this.minDeliveryDate = json.min_delivery_date as string;
		this.maxDeliveryDate = json.max_delivery_date as string;

		this.shipments = [];
		for (const shipmentJSON of json.shipments as JSONArray) {
			const shipment = new Shipment();
			shipment.fromJSON(shipmentJSON as JSONObject);
			this.shipments.push(shipment);
		}
	}

	toJSON() {
		return {
			shipping: this.shipping,
			shipping_method_name: this.shippingMethodName,
			rate: this.rate,
			currency: this.currency,
			min_delivery_days: this.minDeliveryDays,
			max_delivery_days: this.maxDeliveryDays,
			min_delivery_date: this.minDeliveryDate,
			max_delivery_date: this.maxDeliveryDate,
			shipments: this.shipments.map(shipment => shipment.toJSON()),
		}
	}
}

export class Shipment {
	departureCountry = '';
	shipmentItems: ShipmentItem[] = [];

	fromJSON(json: JSONObject): void {
		this.departureCountry = json.departure_country as string;

		this.shipmentItems = [];
		for (const shipmentItemsJSON of json.shipment_items as JSONArray) {
			const shipmentItem = new ShipmentItem();
			shipmentItem.fromJSON(shipmentItemsJSON as JSONObject);
			this.shipmentItems.push(shipmentItem);
		}
	}

	toJSON() {
		return {
			departure_country: this.departureCountry,
			shipment_items: this.shipmentItems.map(shipmentItem => shipmentItem.toJSON()),
		}
	}
}

export class ShipmentItem {
	catalogVariantId = 0;
	quantity = 0;

	fromJSON(json: JSONObject): void {
		this.catalogVariantId = Number(json.catalog_variant_id);
		this.quantity = Number(json.quantity);
	}

	toJSON() {
		return {
			catalog_variant_id: this.catalogVariantId,
			quantity: this.quantity,
		}
	}
}
