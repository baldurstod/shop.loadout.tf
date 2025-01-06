import { createElement, createShadowRoot } from 'harmony-ui';
import productPageCSS from '../../css/productpage.css';
import { Product } from '../model/product';
import { defineShopProduct, HTMLShopProductElement } from './components/shopproduct';
import { defineColumnCart } from './components/columncart';
import { ShopElement } from './shopelement';

export class ProductPage extends ShopElement {
	#htmlShopProduct?: HTMLShopProductElement;
	#htmlColumnCart?: HTMLElement;

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		defineShopProduct();
		defineColumnCart();
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: productPageCSS,
			childs: [
				this.#htmlShopProduct = createElement('shop-product') as HTMLShopProductElement,
				this.#htmlColumnCart = createElement('column-cart'),
			],
		});
	}

	setProduct(product: Product) {
		this.initHTML();
		this.#htmlShopProduct!.setProduct(product);
	}
}
