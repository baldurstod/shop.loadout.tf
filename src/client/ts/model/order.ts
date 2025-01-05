import { Address } from './address';
import { OrderItem } from './orderitem';
import { ShippingInfo } from './shippinginfo';
import { TaxInfo } from './taxinfo';
import { DEFAULT_SHIPPING_METHOD } from '../constants';
import { roundPrice } from '../common';
import { JSONArray, JSONObject } from '../types';

export class Order {
	#id: string = '';
	#currency: string = 'USD';
	#creationTime: number = 0;
	#shippingAddress = new Address();
	#billingAddress = new Address();
	#sameBillingAddress: boolean = true;
	#items: Array<OrderItem> = [];
	#shippingInfos = new Map<string, ShippingInfo>();
	#taxInfo = new TaxInfo();
	#shippingMethod = DEFAULT_SHIPPING_METHOD;
	#printfulOrder: any/*TODO: improve type*/;
	#paypalOrderId: any/*TODO: improve type*/;
	#status = 'created';

	set items(items) {
		this.#items.length = 0;
		if (items) {
			for (let item of items) {
				this.#items.push(item);
			}
		}
	}

	get id() {
		return this.#id;
	}

	get items() {
		return this.#items;
	}

	get shippingAddress() {
		return this.#shippingAddress;
	}

	get billingAddress() {
		return this.#sameBillingAddress ? this.#shippingAddress : this.#billingAddress;
	}

	get sameBillingAddress() {
		return this.#sameBillingAddress;
	}

	getSameBillingAddress() {
		return this.#sameBillingAddress;
	}

	set sameBillingAddress(sameBillingAddress) {
		this.#sameBillingAddress = sameBillingAddress;
	}

	setSameBillingAddress(sameBillingAddress: boolean) {
		this.#sameBillingAddress = sameBillingAddress;
	}

	addShippingInfo(shippingInfo: ShippingInfo) {
		this.#shippingInfos.set(shippingInfo.id, shippingInfo);
	}

	get shippingInfos() {
		return this.#shippingInfos;
	}

	get shippingInfo() {
		return this.#shippingInfos.get(this.#shippingMethod) ?? this.#shippingInfos.get(DEFAULT_SHIPPING_METHOD);
	}

	set taxInfo(taxInfo) {
		this.#taxInfo = taxInfo;
	}

	get taxInfo() {
		return this.#taxInfo;
	}

	set currency(currency) {
		this.#currency = currency;
	}

	get currency() {
		return this.#currency;
	}

	set creationTime(creationTime) {
		this.#creationTime = creationTime;
	}

	get creationTime() {
		return this.#creationTime;
	}

	set paypalOrderId(paypalOrderId) {
		this.#paypalOrderId = paypalOrderId;
	}

	get paypalOrderId() {
		return this.#paypalOrderId;
	}

	set status(status) {
		this.#status = status;
	}

	get status() {
		return this.#status;
	}

	get itemsPrice() {
		let price = 0;
		for (let item of this.#items) {
			price += item.getSubtotal();
		}
		return this.roundPrice(price);
	}

	get shippingPrice() {
		if (this.shippingInfo) {
			return this.roundPrice(this.shippingInfo.rate);
		}
		return 0;
	}

	get taxPrice() {
		if (this.shippingInfo && this.#taxInfo) {
			return this.roundPrice(this.itemsPrice * this.#taxInfo.rate + this.shippingInfo.rate * this.#taxInfo.rate * (this.#taxInfo.shippingTaxable ? 1 : 0));
		}
		return 0;
	}

	get totalPrice() {
		if (this.shippingInfo && this.#taxInfo) {
			return this.itemsPrice + this.shippingPrice + this.taxPrice;
		}
	}

	get shippingMethod() {
		return this.#shippingMethod;
	}

	set shippingMethod(shippingMethod) {
		this.#shippingMethod = shippingMethod;
	}

	get externalId() {
		return this.#printfulOrder?.external_id;
	}

	roundPrice(price: number) {
		return roundPrice(this.#currency, price);
	}

	fromJSON(json: JSONObject) {
		this.#id = json.id as string;
		this.#currency = json.currency as string;
		this.#creationTime = json.date_created as number;
		this.#shippingAddress.fromJSON(json.shipping_address as JSONObject);
		this.#billingAddress.fromJSON(json.billing_address as JSONObject);
		this.#sameBillingAddress = json.same_billing_address as boolean;
		this.#items = [];
		if (json.items) {
			(json.items as JSONArray).forEach(
				(item) => {
					let orderItem = new OrderItem();
					orderItem.fromJSON(item as JSONObject);
					this.#items.push(orderItem);
				}
			);
		}

		this.#shippingInfos.clear();
		const shippingInfos = json.shipping_infos as JSONObject;
		if (shippingInfos) {
			for (const shippingMethod in shippingInfos) {
				let shippingInfo = new ShippingInfo();
				shippingInfo.fromJSON(shippingInfos[shippingMethod] as JSONObject);
				this.addShippingInfo(shippingInfo);
			}
		}

		this.#taxInfo.fromJSON(json.tax_info as JSONObject);
		this.#shippingMethod = json.shipping_method as string;
		this.#printfulOrder = json.printfulOrder;
		this.#paypalOrderId = json.paypalOrderId;
		this.#status = json.status as string;
	}

	toJSON() {
		const itemsJSON: JSONArray = [];
		this.#items.forEach(item => itemsJSON.push(item.toJSON()));

		return {
			id: this.id,
			currency: this.currency,
			creationTime: this.creationTime,
			shipping_address: this.shippingAddress.toJSON(),
			billing_address: this.billingAddress.toJSON(),
			same_billing_address: this.sameBillingAddress,
			items: itemsJSON,
			shipping_infos: this.#shippingInfos,
			taxInfo: this.taxInfo.toJSON(),
			shippingMethod: this.shippingMethod,
			printfulOrder: this.#printfulOrder,
			paypalOrderId: this.paypalOrderId,
			status: this.status,
			priceBreakDown: {
				itemsPrice: this.itemsPrice,
				shippingPrice: this.shippingPrice,
				taxPrice: this.taxPrice,
				totalPrice: this.totalPrice,
			}
		}
	}
}
