import { I18n, createElement, shadowRootStyle } from 'harmony-ui';
import orderSummaryCSS from '../../../css/ordersummary.css';
import { Controller, ControllerEvent, NavigateToDetail } from '../../controller';
import { Order } from '../../model/order';
import { OrderItem } from '../../model/orderitem';
import { formatPercent, formatPrice } from '../../utils';

export class HTMLOrderSummaryElement extends HTMLElement {
	#shadowRoot!: ShadowRoot;
	#htmlProducts?: HTMLElement;
	#htmlPriceBreakDown?: HTMLElement;
	#htmlItemsPrice?: HTMLElement;
	#htmlTaxRate?: HTMLElement;
	#htmlTaxPrice?: HTMLElement;
	#htmlShippingPrice?: HTMLElement;
	#htmlTotalPrice?: HTMLElement;
	#order?: Order;

	constructor() {
		super();
		this.#initHTML();
	}

	#initHTML(): void {
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
		shadowRootStyle(this.#shadowRoot, orderSummaryCSS);

		createElement('div', {
			class: 'summary',
			parent: this.#shadowRoot,
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
		});
	}

	#refreshHTML(): void {
		if (!this.#order) {
			return;
		}
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

	#htmlItemSummary(item: OrderItem, currency: string): HTMLElement {
		const url = `/@product/${item.getId()}`;
		return createElement('div', {
			class: 'item-summary',
			childs: [
				createElement('img', { class: 'thumb', src: item.getThumbnailUrl() }),
				createElement('div', { class: 'name', innerText: item.getName() }),
				createElement('div', { class: 'unit-price', innerText: formatPrice(item.getRetailPrice(), currency) }),
				createElement('div', { class: 'quantity', innerText: String(item.getQuantity()) }),
				createElement('div', { class: 'total-price', innerText: formatPrice(item.getRetailPrice() * item.getQuantity(), currency) }),
			],
			$click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url } }),
			$mouseup: (event: MouseEvent) => {
				if (event.button == 1) {
					open(url, '_blank');
				}
			},
		});
	}

	setOrder(order: Order): void {
		this.#order = order;
		this.#refreshHTML();
	}
}

let definedOrderSummary = false;
export function defineOrderSummary(): void {
	if (window.customElements && !definedOrderSummary) {
		customElements.define('order-summary', HTMLOrderSummaryElement);
		definedOrderSummary = true;
	}
}
