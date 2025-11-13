import { createElement, createShadowRoot } from 'harmony-ui';
import { Order } from '../model/order';
import { OrderItem } from '../model/orderitem';
import { formatPercent, formatPrice } from '../utils';
import { ShopElement } from './shopelement';

export class OrderSummary extends ShopElement {
	#htmlProducts?: HTMLElement;
	#htmlSubtotal?: HTMLElement;
	#htmlShippingPrice?: HTMLElement;
	#htmlTaxRate?: HTMLElement;
	#htmlTax?: HTMLElement;
	#htmlTotal?: HTMLElement;

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('div');

		this.#htmlProducts = createElement('div', { class: 'order-summary-products' });
		const htmlSubTotalContainer = createElement('div', { class: 'order-summary-total-container' });
		const htmlTotalContainer = createElement('div', { class: 'order-summary-total-container' });

		const htmlSubtotalLine = createElement('div');
		this.#htmlSubtotal = createElement('span', { class: 'order-summary-subtotal' });
		htmlSubtotalLine.append(createElement('label', { i18n: '#subtotal' }), this.#htmlSubtotal);

		const htmlShippingLine = createElement('div');
		this.#htmlShippingPrice = createElement('span');
		htmlShippingLine.append(createElement('label', { i18n: '#shipping' }), this.#htmlShippingPrice);

		const htmlTaxLine = createElement('div');
		this.#htmlTaxRate = createElement('span');
		this.#htmlTax = createElement('span');
		htmlTaxLine.append(createElement('label', { childs: [createElement('span', { i18n: '#tax' }), this.#htmlTaxRate] }), this.#htmlTax);

		const htmlTotalLine = createElement('div');
		this.#htmlTotal = createElement('span', { class: 'order-summary-total' });
		htmlTotalLine.append(createElement('label', { i18n: '#total' }), this.#htmlTotal);

		htmlSubTotalContainer.append(htmlSubtotalLine, htmlShippingLine, htmlTaxLine);
		htmlTotalContainer.append(htmlTotalLine);
		this.shadowRoot.append(this.#htmlProducts, htmlSubTotalContainer, htmlTotalContainer);
	}

	#refreshHTML(order: Order | null): void {
		this.initHTML();

		//this.htmlElement.innerText = '';
		this.#htmlProducts!.innerText = '';
		this.#htmlShippingPrice!.innerText = '';
		this.#htmlTaxRate!.innerText = '';
		this.#htmlTax!.innerText = '';
		this.#htmlTotal!.innerText = '';

		if (order) {
			order.items.forEach((item) => {
				this.#htmlProducts!.append(this.#htmlItemSummary(item, order.currency));
			});

			this.#htmlSubtotal!.innerText = formatPrice(order.itemsPrice, order.currency);



			if (order.taxInfo) {
				this.#htmlTaxRate!.innerText = ` (${formatPercent(order.taxInfo.rate)})`;
			}

			if (order.taxPrice) {
				this.#htmlTax!.innerText = formatPrice(order.taxPrice, order.currency);
			}

			if (order.shippingPrice) {
				this.#htmlShippingPrice!.innerText = formatPrice(order.shippingPrice, order.currency);
			}

			if (order.totalPrice) {
				this.#htmlTotal!.innerText = formatPrice(order.totalPrice, order.currency);
			}
		}

		/*if (summary.shippingInfo) {
			this.htmlShipping.innerText = formatPrice(summary.shippingInfo.rate, summary.shippingInfo.currency);
			if (summary.taxInfo) {
				this.htmlTaxRate.innerText = ` (${formatPercent(summary.taxInfo.rate)})`;
				this.htmlTax.innerText = formatPrice(summary.taxPrice, summary.currency);
				this.htmlTotal.innerText = formatPrice(summary.totalPrice, summary.currency);
			}
		} else {
			this.htmlShipping.append(createElement('label', {i18n:'#calculated_at_next_step'}));
		}*/

	}


	#htmlItemSummary(item: OrderItem, currency: string): HTMLElement {
		const htmlSummary = createElement('div', { class: 'item-summary' });
		const htmlProductThumb = createElement('img', { class: 'thumb', src: item.getThumbnailUrl() });
		const htmlProductName = createElement('div', { class: 'name', innerText: item.getName() });
		const htmlProductPrice = createElement('div', { class: 'price', innerText: formatPrice(item.getRetailPrice(), currency) });
		const htmlProductQuantity = createElement('div', { class: 'quantity', innerText: String(item.getQuantity()) });

		htmlSummary.append(htmlProductThumb, htmlProductQuantity, htmlProductName, htmlProductPrice);
		return htmlSummary;
	}

	setOrder(order: Order | null): void {
		this.#refreshHTML(order);
	}
}
