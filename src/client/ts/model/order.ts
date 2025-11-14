import { roundPrice } from '../common';
import { DEFAULT_SHIPPING_METHOD } from '../constants';
import { OrderItemJSON, OrderJSON, ShippingRateJSON } from '../responses/order';
import { Address } from './address';
import { OrderItem } from './orderitem';
import { ShippingInfo } from './shippinginfo';
import { TaxInfo } from './taxinfo';

export class Order {
	#id = '';
	#currency = 'USD';
	#dateCreated = 0;
	#dateUpdated = 0;
	readonly #shippingAddress = new Address();
	readonly #billingAddress = new Address();
	#sameBillingAddress = true;
	#items: OrderItem[] = [];
	#shippingInfos = new Map<string, ShippingInfo>();
	#taxInfo = new TaxInfo();
	#shippingMethod = DEFAULT_SHIPPING_METHOD;
	#printfulOrderId = '';
	#paypalOrderId = '';
	#status = 'created';

	set items(items) {
		this.#items.length = 0;
		if (items) {
			for (const item of items) {
				this.#items.push(item);
			}
		}
	}

	get items(): OrderItem[] {
		return this.#items;
	}

	get id(): string {
		return this.#id;
	}

	get shippingAddress(): Address {
		return this.#shippingAddress;
	}

	get billingAddress(): Address {
		return this.#sameBillingAddress ? this.#shippingAddress : this.#billingAddress;
	}

	get sameBillingAddress(): boolean {
		return this.#sameBillingAddress;
	}

	set sameBillingAddress(sameBillingAddress) {
		this.#sameBillingAddress = sameBillingAddress;
	}

	getSameBillingAddress(): boolean {
		return this.#sameBillingAddress;
	}

	setSameBillingAddress(sameBillingAddress: boolean): void {
		this.#sameBillingAddress = sameBillingAddress;
	}

	addShippingInfo(shippingInfo: ShippingInfo): void {
		this.#shippingInfos.set(shippingInfo.shipping, shippingInfo);
	}

	get shippingInfos(): Map<string, ShippingInfo> {
		return this.#shippingInfos;
	}

	get shippingInfo(): ShippingInfo | null {
		return this.#shippingInfos.get(this.#shippingMethod) ?? this.#shippingInfos.get(DEFAULT_SHIPPING_METHOD) ?? null;
	}

	set taxInfo(taxInfo) {
		this.#taxInfo = taxInfo;
	}

	get taxInfo(): TaxInfo {
		return this.#taxInfo;
	}

	set currency(currency) {
		this.#currency = currency;
	}

	get currency(): string {
		return this.#currency;
	}

	setDateCreated(dateCreated: number): void {
		this.#dateCreated = dateCreated;
	}

	getDateCreated(): number {
		return this.#dateCreated;
	}

	setUpdated(dateUpdated: number): void {
		this.#dateUpdated = dateUpdated;
	}

	getUpdated(): number {
		return this.#dateUpdated;
	}

	set paypalOrderId(paypalOrderId) {
		this.#paypalOrderId = paypalOrderId;
	}

	get paypalOrderId(): string {
		return this.#paypalOrderId;
	}

	set status(status) {
		this.#status = status;
	}

	get status(): string {
		return this.#status;
	}

	get itemsPrice(): number {
		let price = 0;
		for (const item of this.#items) {
			price += item.getSubtotal();
		}
		return this.roundPrice(price);
	}

	get shippingPrice(): number {
		if (this.shippingInfo) {
			return this.roundPrice(Number(this.shippingInfo.rate));
		}
		return 0;
	}

	get taxPrice(): number {
		if (this.shippingInfo && this.#taxInfo) {
			return this.roundPrice(this.itemsPrice * this.#taxInfo.rate + Number(this.shippingInfo.rate) * this.#taxInfo.rate * (this.#taxInfo.shippingTaxable ? 1 : 0));
		}
		return 0;
	}

	get totalPrice(): number | undefined {
		if (this.shippingInfo && this.#taxInfo) {
			return this.itemsPrice + this.shippingPrice + this.taxPrice;
		}
	}

	get shippingMethod(): string {
		return this.#shippingMethod;
	}

	set shippingMethod(shippingMethod) {
		this.#shippingMethod = shippingMethod;
	}

	get externalId(): string {
		return this.#printfulOrderId;
	}

	roundPrice(price: number): number {
		return roundPrice(this.#currency, price);
	}

	fromJSON(json: OrderJSON): void {
		this.#id = json.id;
		this.#currency = json.currency;
		this.#dateCreated = json.date_created;
		this.#dateUpdated = json.date_updated;
		this.#shippingAddress.fromJSON(json.shipping_address);
		this.#billingAddress.fromJSON(json.billing_address);
		this.#sameBillingAddress = json.same_billing_address;
		this.#items = [];
		for (const item of json.items) {
			const orderItem = new OrderItem();
			orderItem.fromJSON(item);
			this.#items.push(orderItem);
		}

		this.#shippingInfos.clear();
		const shippingInfos = json.shipping_infos;
		if (shippingInfos) {
			for (const shippingInfoJSON of shippingInfos) {
				const shippingInfo = new ShippingInfo();
				shippingInfo.fromJSON(shippingInfoJSON);
				this.addShippingInfo(shippingInfo);
			}
		}

		this.#taxInfo.fromJSON(json.tax_info);
		this.#shippingMethod = json.shipping_method;
		this.#printfulOrderId = json.printful_order_id;
		this.#paypalOrderId = json.paypal_order_id;
		this.#status = json.status;
	}

	toJSON(): OrderJSON {
		const itemsJSON: OrderItemJSON[] = [];
		this.#items.forEach(item => itemsJSON.push(item.toJSON()));

		const shippingInfosJSON: ShippingRateJSON[] = [];
		for (const [, shippingInfos] of this.#shippingInfos) {

			shippingInfosJSON.push(shippingInfos.toJSON());
		}

		return {
			id: this.id,
			currency: this.currency,
			date_created: this.#dateCreated,
			date_updated: this.#dateUpdated,
			shipping_address: this.shippingAddress.toJSON(),
			billing_address: this.billingAddress.toJSON(),
			same_billing_address: this.sameBillingAddress,
			items: itemsJSON,
			shipping_infos: shippingInfosJSON,
			tax_info: this.taxInfo.toJSON(),
			shipping_method: this.shippingMethod,
			printful_order_id: this.#printfulOrderId,
			paypal_order_id: this.paypalOrderId,
			status: this.status,
			/*
			priceBreakDown: {
				itemsPrice: this.itemsPrice,
				shippingPrice: this.shippingPrice,
				taxPrice: this.taxPrice,
				totalPrice: this.totalPrice,
			}
			*/
		};
	}
}
