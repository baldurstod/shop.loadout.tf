import { createElement } from 'harmony-ui';
import productPageCSS from '../../css/productpage.css';
import { Product } from '../model/product';
import { defineShopProduct, HTMLShopProductElement } from './components/shopproduct';
import { defineColumnCart } from './components/columncart';

export class ProductPage {
	#htmlElement: HTMLElement;
	#htmlShopProduct: HTMLShopProductElement;
	#htmlColumnCart: HTMLElement;

	#initHTML() {
		defineShopProduct();
		defineColumnCart();
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: productPageCSS,
			childs: [
				this.#htmlShopProduct = createElement('shop-product') as HTMLShopProductElement,
				this.#htmlColumnCart = createElement('column-cart'),
			],
		});
		return this.#htmlElement;
	}

	get htmlElement(): HTMLElement {
		return this.#htmlElement ?? this.#initHTML();
	}

	setProduct(product: Product) {
		this.#htmlShopProduct.setProduct(product);
	}
}
