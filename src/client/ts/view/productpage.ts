import { createElement } from 'harmony-ui';
export { ShopProduct } from './components/shopproduct.js';
export { ColumnCart } from './components/columncart.js';

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
		return;
		for (const shopProduct of products) {
			createElement('shop-product', {
				parent: this.#htmlElement,
				elementCreated: element => element.setProduct(shopProduct),
			});
		}
	}
}
