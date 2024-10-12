import { createElement } from 'harmony-ui';
export * from './components/columncart';
export * from './components/shopproduct';

import productPageCSS from '../../css/productpage.css';

export class ProductPage {
	#htmlElement;
	#htmlShopProduct;
	#htmlColumnCart;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: productPageCSS,
			childs: [
				this.#htmlShopProduct = createElement('shop-product'),
				this.#htmlColumnCart = createElement('column-cart'),
			],
		});
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setProduct(product) {
		this.#htmlShopProduct.setProduct(product);
	}
}
