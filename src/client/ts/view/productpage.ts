import { createElement, createShadowRoot } from 'harmony-ui';
import productPageCSS from '../../css/productpage.css';
import { Product } from '../model/product';
import { defineShopProduct, HTMLShopProductElement } from './components/shopproduct';
import { defineColumnCart } from './components/columncart';

export class ProductPage {
	#shadowRoot?: ShadowRoot;
	#htmlShopProduct?: HTMLShopProductElement;
	#htmlColumnCart?: HTMLElement;

	#initHTML() {
		defineShopProduct();
		defineColumnCart();
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyle: productPageCSS,
			childs: [
				this.#htmlShopProduct = createElement('shop-product') as HTMLShopProductElement,
				this.#htmlColumnCart = createElement('column-cart'),
			],
		});
		return this.#shadowRoot.host;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}

	setProduct(product: Product) {
		if (!this.#htmlShopProduct) {
			this.#initHTML();
		}
		this.#htmlShopProduct!.setProduct(product);
	}
}
