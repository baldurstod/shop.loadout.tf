import { createElement } from 'harmony-ui';
export { ShopProductWidgetElement } from './components/shopproductwidget.js';

import productsPageCSS from '../../css/productspage.css';

export class ProductsPage {
	#htmlElement;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: productsPageCSS,
			childs: [
			],
		});
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setProducts(products) {
		for (const shopProduct of products) {
			createElement('shop-product-widget', {
				parent: this.#htmlElement,
				elementCreated: element => element.setProduct(shopProduct),
			});
		}
	}
}
