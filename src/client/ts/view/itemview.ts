import { formatPrice } from '../utils';
import { createElement } from 'harmony-ui';

export class ItemView extends EventTarget {
	#model;
	constructor(model) {
		super();
		this.model = model;
		this.#initHTML();
	}

	set model(model) {
		if (this.#model != model) {
			this.#model = model;
			model.addEventListener('change', () => this.#refreshHTML());
		}
	}

	get html() {
		this.#refreshHTML();
		return this.htmlElement;
	}

	htmlSummary(currency) {
		let htmlSummary = createElement('div', { class: 'item-summary' });
		let htmlProductThumb = createElement('img', { class: 'thumb', src: this.#model.getFileUrl('thumbnail')?.url });
		let htmlProductName = createElement('div', { class: 'name', innerHTML: this.#model.name });
		let htmlProductPrice = createElement('td', { class: 'price', innerHTML: formatPrice(this.#model.retailPrice, currency) });
		let htmlProductQuantity = createElement('div', { class: 'quantity', innerHTML: this.#model.quantity });
		/*

			let htmlElement = createElement('div', {class:'order-summary-product'});
			let htmlProductThumb = createElement('img', {class:'thumb',src:this.#thumbnailUrl});
			let htmlProductName = createElement('div', {class:'name',innerHTML:this.#name});
			let htmlProductPrice = createElement('td', {class:'price',innerHTML:formatPrice(this.#price, currency)});
			let htmlProductQuantity = createElement('div', {class:'quantity',innerHTML:this.#quantity});

			htmlElement.append(htmlProductThumb, htmlProductQuantity, htmlProductName, htmlProductPrice);
			return htmlElement;*/





		htmlSummary.append(htmlProductThumb, htmlProductQuantity, htmlProductName, htmlProductPrice);
		return htmlSummary;
	}


	#initHTML() {
		this.htmlElement = createElement('div', { class: 'item' });


		this.htmlElement.append();
	}

	#refreshHTML() {

	}
}
