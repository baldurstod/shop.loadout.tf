import { createElement } from 'harmony-ui';
export * from './components/columncart';
export * from './components/shopproduct';
import productPageCSS from '../../css/productpage.css';
import { Product } from '../model/product';
import { ShopProductElement } from './productpage';

export class ProductPage {
	#htmlElement: HTMLElement;
	#htmlShopProduct: ShopProductElement;
	#htmlColumnCart: HTMLElement;

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

	get htmlElement(): HTMLElement {
		return this.#htmlElement ?? this.#initHTML();
	}

	setProduct(product: Product) {
		this.#htmlShopProduct.setProduct(product);
	}
}
