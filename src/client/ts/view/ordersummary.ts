import { createElement, createShadowRoot } from 'harmony-ui';
import { formatPrice, formatPercent } from '../utils';
import { Order } from '../model/order';
import { OrderItem } from '../model/orderitem.js';

export class OrderSummary {
	#shadowRoot?: ShadowRoot;
	#htmlProducts?: HTMLElement;
	#htmlSubtotal?: HTMLElement;
	#htmlShippingPrice?: HTMLElement;
	#htmlTaxRate?: HTMLElement;
	#htmlTax?: HTMLElement;
	#htmlTotal?: HTMLElement;

	#initHTML() {
		this.#shadowRoot = createShadowRoot('div');

		this.#htmlProducts = createElement('div', { class: 'order-summary-products' });
		let htmlSubTotalContainer = createElement('div', { class: 'order-summary-total-container' });
		let htmlTotalContainer = createElement('div', { class: 'order-summary-total-container' });

		let htmlSubtotalLine = createElement('div');
		this.#htmlSubtotal = createElement('span', { class: 'order-summary-subtotal' });
		htmlSubtotalLine.append(createElement('label', { i18n: '#subtotal' }), this.#htmlSubtotal);

		let htmlShippingLine = createElement('div');
		this.#htmlShippingPrice = createElement('span');
		htmlShippingLine.append(createElement('label', { i18n: '#shipping' }), this.#htmlShippingPrice);

		let htmlTaxLine = createElement('div');
		this.#htmlTaxRate = createElement('span');
		this.#htmlTax = createElement('span');
		htmlTaxLine.append(createElement('label', { childs: [createElement('span', { i18n: '#tax' }), this.#htmlTaxRate] }), this.#htmlTax);

		let htmlTotalLine = createElement('div');
		this.#htmlTotal = createElement('span', { class: 'order-summary-total' });
		htmlTotalLine.append(createElement('label', { i18n: '#total' }), this.#htmlTotal);

		htmlSubTotalContainer.append(htmlSubtotalLine, htmlShippingLine, htmlTaxLine);
		htmlTotalContainer.append(htmlTotalLine);
		this.#shadowRoot.append(this.#htmlProducts, htmlSubTotalContainer, htmlTotalContainer);
	}

	#refreshHTML(order: Order | null) {
		if (!this.#shadowRoot) {
			this.#initHTML();
		}

		//this.htmlElement.innerHTML = '';
		this.#htmlProducts!.innerText = '';
		this.#htmlShippingPrice!.innerText = '';
		this.#htmlTaxRate!.innerText = '';
		this.#htmlTax!.innerText = '';
		this.#htmlTotal!.innerText = '';

		if (order) {
			order.items.forEach((item) => {
				this.#htmlProducts!.append(this.#htmlItemSummary(item, order.currency));
			});

			this.#htmlSubtotal!.innerHTML = formatPrice(order.itemsPrice, order.currency);



			if (order.taxInfo) {
				this.#htmlTaxRate!.innerHTML = ` (${formatPercent(order.taxInfo.rate)})`;
			}

			if (order.taxPrice) {
				this.#htmlTax!.innerHTML = formatPrice(order.taxPrice, order.currency);
			}

			if (order.shippingPrice) {
				this.#htmlShippingPrice!.innerHTML = formatPrice(order.shippingPrice, order.currency);
			}

			if (order.totalPrice) {
				this.#htmlTotal!.innerHTML = formatPrice(order.totalPrice, order.currency);
			}
		}

		/*if (summary.shippingInfo) {
			this.htmlShipping.innerHTML = formatPrice(summary.shippingInfo.rate, summary.shippingInfo.currency);
			if (summary.taxInfo) {
				this.htmlTaxRate.innerHTML = ` (${formatPercent(summary.taxInfo.rate)})`;
				this.htmlTax.innerHTML = formatPrice(summary.taxPrice, summary.currency);
				this.htmlTotal.innerHTML = formatPrice(summary.totalPrice, summary.currency);
			}
		} else {
			this.htmlShipping.append(createElement('label', {i18n:'#calculated_at_next_step'}));
		}*/

	}


	#htmlItemSummary(item: OrderItem, currency: string) {
		let htmlSummary = createElement('div', { class: 'item-summary' });
		let htmlProductThumb = createElement('img', { class: 'thumb', src: item.getThumbnailUrl() });
		let htmlProductName = createElement('div', { class: 'name', innerText: item.getName() });
		let htmlProductPrice = createElement('div', { class: 'price', innerText: formatPrice(item.getRetailPrice(), currency) });
		let htmlProductQuantity = createElement('div', { class: 'quantity', innerText: item.getQuantity() });

		htmlSummary.append(htmlProductThumb, htmlProductQuantity, htmlProductName, htmlProductPrice);
		return htmlSummary;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}

	setOrder(order: Order | null) {
		this.#refreshHTML(order);
	}
}
