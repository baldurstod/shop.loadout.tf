import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import orderPageCSS from '../../css/orderpage.css';
import { Order } from '../model/order';
import { ShopElement } from './shopelement';
import { OrderItem } from '../model/orderitem';
import { formatPercent, formatPrice } from '../utils';

export class OrderPage extends ShopElement {
	#order?: Order;
	#htmlProducts?: HTMLElement;
	#htmlPriceBreakDown?: HTMLElement;
	#htmlItemsPrice?: HTMLElement;
	#htmlTaxRate?: HTMLElement;
	#htmlTaxPrice?: HTMLElement;
	#htmlShippingPrice?: HTMLElement;
	#htmlTotalPrice?: HTMLElement;

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: orderPageCSS,
			childs: [
				createElement('div', { class: 'title', i18n: '#order_summary' }),
				createElement('div', {
					class: 'summary',
					childs: [
						this.#htmlProducts = createElement('div', { class: 'items' }),
						this.#htmlPriceBreakDown = createElement('div', {
							class: 'price-breakdown',
							childs: [
								createElement('div', {
									class: 'price-line',
									childs: [
										createElement('div', { class: "price-label", i18n: "#subtotal" }),
										this.#htmlItemsPrice = createElement('div', { class: "price" }),
									]
								}),
								createElement('div', {
									class: 'price-line',
									childs: [
										createElement('div', { class: "price-label", i18n: "#shipping" }),
										this.#htmlShippingPrice = createElement('div', { class: "price" }),
									]
								}),
								createElement('div', {
									class: 'price-line',
									childs: [
										createElement('div', { class: "price-label", i18n: "#tax" }),
										this.#htmlTaxRate = createElement('div'),
										this.#htmlTaxPrice = createElement('div', { class: "price" }),
									]
								}),
								createElement('div', {
									class: 'price-line',
									childs: [
										createElement('div', { class: "price-label", i18n: "#total" }),
										this.#htmlTotalPrice = createElement('div', { class: "price" }),
									]
								}),
							],
						}),
					],
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#refreshHTML() {
		if (!this.#order) {
			return;
		}
		this.initHTML();
		this.#htmlProducts!.innerText = '';

		this.#order.items.forEach((item) => {
			this.#htmlProducts!.append(this.#htmlItemSummary(item, this.#order!.currency));
		});


		this.#htmlItemsPrice!.innerText = formatPrice(this.#order.itemsPrice, this.#order.currency);

		this.#htmlTaxRate!.innerText = '';
		if (this.#order.taxInfo) {
			this.#htmlTaxRate!.innerText = ` (${formatPercent(this.#order.taxInfo.rate)})`;
		}

		if (this.#order.taxPrice) {
			this.#htmlTaxPrice!.innerText = formatPrice(this.#order.taxPrice, this.#order.currency);
		}

		if (this.#order.shippingPrice) {
			this.#htmlShippingPrice!.innerText = formatPrice(this.#order.shippingPrice, this.#order.currency);
		}

		if (this.#order.totalPrice) {
			this.#htmlTotalPrice!.innerText = formatPrice(this.#order.totalPrice, this.#order.currency);
		}

	}

	#htmlItemSummary(item: OrderItem, currency: string) {
		return createElement('div', {
			class: 'item-summary',
			childs: [
				createElement('img', { class: 'thumb', src: item.getThumbnailUrl() }),
				createElement('div', { class: 'name', innerText: item.getName() }),
				createElement('div', { class: 'quantity', innerText: item.getQuantity() }),
				createElement('div', { class: 'price', innerText: formatPrice(item.getRetailPrice(), currency) }),
			]
		});
	}

	setOrder(order: Order) {
		this.#order = order;
		this.#refreshHTML();
	}
}
