import { createElement } from 'harmony-ui';

import productsPageCSS from '../../css/productspage.css';
import { defineShopProductWidget } from './components/shopproductwidget';

export class ProductsPage {
	#htmlElement;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: productsPageCSS,
		});
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setProducts(products = []) {
		defineShopProductWidget();
		this.#htmlElement.innerHTML = '';
		for (const shopProduct of products) {
			createElement('shop-product-widget', {
				parent: this.#htmlElement,
				elementCreated: element => element.setProduct(shopProduct),
			});
		}
	}
}
