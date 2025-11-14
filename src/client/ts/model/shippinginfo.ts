import { ShipmentItemJSON, ShipmentsJSON, ShippingRateJSON } from '../responses/order';

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

	fromJSON(json: ShippingRateJSON): void {
		this.shipping = json.shipping;
		this.shippingMethodName = json.shipping_method_name;
		this.rate = json.rate;
		this.currency = json.currency;
		this.minDeliveryDays = json.min_delivery_days;
		this.maxDeliveryDays = json.max_delivery_days;
		this.minDeliveryDate = json.min_delivery_date;
		this.maxDeliveryDate = json.max_delivery_date;

		this.shipments = [];
		for (const shipmentJSON of json.shipments) {
			const shipment = new Shipment();
			shipment.fromJSON(shipmentJSON);
			this.shipments.push(shipment);
		}
	}

	toJSON(): ShippingRateJSON {
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
	customsFeesPossible = false;

	fromJSON(json: ShipmentsJSON): void {
		this.departureCountry = json.departure_country;

		this.shipmentItems = [];
		for (const shipmentItemsJSON of json.shipment_items) {
			const shipmentItem = new ShipmentItem();
			shipmentItem.fromJSON(shipmentItemsJSON);
			this.shipmentItems.push(shipmentItem);
		}

		this.customsFeesPossible = json.customs_fees_possible;
	}

	toJSON(): ShipmentsJSON {
		return {
			departure_country: this.departureCountry,
			shipment_items: this.shipmentItems.map(shipmentItem => shipmentItem.toJSON()),
			customs_fees_possible: this.customsFeesPossible,
		}
	}
}

export class ShipmentItem {
	catalogVariantId = 0;
	quantity = 0;

	fromJSON(json: ShipmentItemJSON): void {
		this.catalogVariantId = Number(json.catalog_variant_id);
		this.quantity = Number(json.quantity);
	}

	toJSON(): ShipmentItemJSON {
		return {
			catalog_variant_id: this.catalogVariantId,
			quantity: this.quantity,
		}
	}
}
