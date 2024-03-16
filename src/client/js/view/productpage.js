import { createElement } from 'harmony-ui';
export { ShopProduct } from './components/shopproduct.js';

import productPageCSS from '../../css/productpage.css';

export class ProductPage {
	#htmlElement;
	#htmlShopProduct;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: productPageCSS,
			childs: [
				this.#htmlShopProduct = createElement('shop-product'),
			],
		});
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setProduct(product) {
		this.#htmlShopProduct.setProduct(product);
		return;
		for (const shopProduct of products) {
			createElement('shop-product', {
				parent: this.#htmlElement,
				elementCreated: element => element.setProduct(shopProduct),
			});
		}
	}
}
