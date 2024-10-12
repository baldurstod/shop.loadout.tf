import { Address } from './address.js';
import { OrderItem } from './orderitem.js';
import { ShippingInfo } from './shippinginfo.js';
import { TaxInfo } from './taxinfo.js';
import { DEFAULT_SHIPPING_METHOD } from '../constants.js';
import { roundPrice } from '../common.js';

export class Order {
	#id;
	#currency;
	#creationTime;
	#shippingAddress;
	#billingAddress;
	#sameBillingAddress;
	#items;
	#shippingInfos = new Map();
	#taxInfo = new TaxInfo();
	#shippingMethod;
	#printfulOrder;
	#paypalOrderId;
	#status = 'created';
	constructor({currency = 'USD', items, id, creationTime}:any = {}) {
		this.#id = id;
		this.#currency = currency;
		this.#creationTime = creationTime;
		this.#shippingAddress = new Address();
		this.#billingAddress = new Address();
		this.#sameBillingAddress = true;
		this.#items = [];
		this.#shippingMethod = DEFAULT_SHIPPING_METHOD;
		//this.#shippingInfo = new ShippingInfo();

/*
id	{…}
external_id	{…}
store	{…}
status	{…}
shipping	{…}
shipping_service_name	{…}
created	{…}
updated	{…}
recipient	{…}
items	{…}
incomplete_items	{…}
costs	{…}
retail_costs	{…}
pricing_breakdown	{…}
shipments	{…}
gift	{…}
packing_slip	{…}*/



		if (items) {
			this.items = items;
		}
	}

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

	setSameBillingAddress(sameBillingAddress) {
		this.#sameBillingAddress = sameBillingAddress;
	}

	addShippingInfo(shippingInfo) {
		this.#shippingInfos.set(shippingInfo.id, shippingInfo);
	}

	get shippingInfos() {
		return this.#shippingInfos;
	}

	get shippingInfo() {
		return this.#shippingInfos.get(this.#shippingMethod) ??  this.#shippingInfos.get(DEFAULT_SHIPPING_METHOD);
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
			price += item.subtotal;
		}
		return this.roundPrice(price);
	}

	get shippingPrice() {
		if (this.shippingInfo) {
			return this.roundPrice(this.shippingInfo.rate);
		}
	}

	get taxPrice() {
		if (this.shippingInfo && this.#taxInfo) {
			return this.roundPrice(this.itemsPrice * this.#taxInfo.rate + this.shippingInfo.rate * this.#taxInfo.rate * (this.#taxInfo.shippingTaxable ? 1 : 0));
		}
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

	roundPrice(price) {
		return roundPrice(this.#currency, price);
	}

	fromJSON(json) {
		this.#id = json.id;
		this.#currency = json.currency;
		this.#creationTime = json.date_created;
		this.#shippingAddress.fromJSON(json.shipping_address);
		this.#billingAddress.fromJSON(json.billing_address);
		this.#sameBillingAddress = json.same_billing_address;
		this.#items = [];
		if (json.items) {
			json.items.forEach(
				(item) => {
					let orderItem = new OrderItem();
					orderItem.fromJSON(item);
					this.#items.push(orderItem);
				}
			);
		}

		this.#shippingInfos.clear();
		let shippingInfos = json.shipping_infos;
		if (shippingInfos) {
			for (const shippingMethod in shippingInfos) {
				let shippingInfo = new ShippingInfo();
				shippingInfo.fromJSON(shippingInfos[shippingMethod]);
				this.addShippingInfo(shippingInfo);
			}
		}

		this.#taxInfo.fromJSON(json.tax_info);
		this.#shippingMethod = json.shipping_method;
		this.#printfulOrder = json.printfulOrder;
		this.#paypalOrderId = json.paypalOrderId;
		this.#status = json.status;
	}

	toJSON() {
		const itemsJSON = [];
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
