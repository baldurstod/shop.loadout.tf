import {createElement} from 'harmony-ui';
import {formatPrice, formatPercent} from '../utils.js';

export class OrderSummary {
	#htmlElement;
	constructor() {
		this.#initHTML();
	}

	#initHTML() {
		this.#htmlElement = createElement('div', {class:'order-summary'});

		this.htmlProducts = createElement('div', {class:'order-summary-products'});
		let htmlSubTotalContainer = createElement('div', {class:'order-summary-total-container'});
		let htmlTotalContainer = createElement('div', {class:'order-summary-total-container'});

		let htmlSubtotalLine = createElement('div');
		this.htmlSubtotal = createElement('span', {class:'order-summary-subtotal'});
		htmlSubtotalLine.append(createElement('label', {i18n:'#subtotal'}), this.htmlSubtotal);

		let htmlShippingLine = createElement('div');
		this.htmlShippingPrice = createElement('span');
		htmlShippingLine.append(createElement('label', {i18n:'#shipping'}), this.htmlShippingPrice);

		let htmlTaxLine = createElement('div');
		this.htmlTaxRate = createElement('span');
		this.htmlTax = createElement('span');
		htmlTaxLine.append(createElement('label', {childs:[createElement('span', {i18n:'#tax'}), this.htmlTaxRate]}), this.htmlTax);

		let htmlTotalLine = createElement('div');
		this.htmlTotal = createElement('span', {class:'order-summary-total'});
		htmlTotalLine.append(createElement('label', {i18n:'#total'}), this.htmlTotal);

		htmlSubTotalContainer.append(htmlSubtotalLine, htmlShippingLine, htmlTaxLine);
		htmlTotalContainer.append(htmlTotalLine);
		this.#htmlElement.append(this.htmlProducts, htmlSubTotalContainer, htmlTotalContainer);
	}

	#refreshHTML(summary) {
		//this.htmlElement.innerHTML = '';
		this.htmlProducts.innerHTML = '';
		this.htmlShippingPrice.innerHTML = '';
		this.htmlTaxRate.innerHTML = '';
		this.htmlTax.innerHTML = '';
		this.htmlTotal.innerHTML = '';

		if (summary) {
			summary.items.forEach((item) => {
				this.htmlProducts.append(this.#htmlItemSummary(item, summary.currency));
			});

			this.htmlSubtotal.innerHTML = formatPrice(summary.itemsPrice, summary.currency);



			if (summary.taxInfo) {
				this.htmlTaxRate.innerHTML = ` (${formatPercent(summary.taxInfo.rate)})`;
			}

			if (summary.taxPrice) {
				this.htmlTax.innerHTML = formatPrice(summary.taxPrice, summary.currency);
			}

			if (summary.shippingPrice) {
				this.htmlShippingPrice.innerHTML = formatPrice(summary.shippingPrice, summary.currency);
			}

			if (summary.totalPrice) {
				this.htmlTotal.innerHTML = formatPrice(summary.totalPrice, summary.currency);
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


	#htmlItemSummary(item, currency) {
		let htmlSummary = createElement('div', {class:'item-summary'});
		let htmlProductThumb = createElement('img', {class:'thumb',src:item.thumbnailUrl});
		let htmlProductName = createElement('div', {class:'name',innerHTML:item.name});
		let htmlProductPrice = createElement('div', {class:'price',innerHTML:formatPrice(item.retailPrice, currency)});
		let htmlProductQuantity = createElement('div', {class:'quantity',innerHTML:item.quantity});

		htmlSummary.append(htmlProductThumb, htmlProductQuantity, htmlProductName, htmlProductPrice);
		return htmlSummary;
	}

	get html() {
		return this.#htmlElement;
	}

	set summary(summary) {
		this.#refreshHTML(summary);
	}
}
