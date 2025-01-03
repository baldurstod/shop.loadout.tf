import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import productsPageCSS from '../../css/productspage.css';
import { defineShopProductWidget, HTMLShopProductWidgetElement } from './components/shopproductwidget';
import { Product } from '../model/product';

export class ProductsPage {
	#shadowRoot?: ShadowRoot;

	#initHTML() {
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyle: productsPageCSS,
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}

	setProducts(products: Array<Product> = []) {
		if (!this.#shadowRoot) {
			this.#initHTML();
		}
		defineShopProductWidget();
		this.#shadowRoot!.innerHTML = '';
		for (const shopProduct of products) {
			createElement('shop-product-widget', {
				parent: this.#shadowRoot,
				elementCreated: (element: HTMLShopProductWidgetElement) => element.setProduct(shopProduct),
			});
		}
	}
}
