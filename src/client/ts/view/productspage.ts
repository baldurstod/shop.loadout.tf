import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import productsPageCSS from '../../css/productspage.css';
import { defineShopProductWidget, HTMLShopProductWidgetElement } from './components/shopproductwidget';
import { Product } from '../model/product';
import { ShopElement } from './shopelement';

export class ProductsPage extends ShopElement {

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: productsPageCSS,
		});
		I18n.observeElement(this.shadowRoot);
	}

	setProducts(products: Array<Product> = []) {
		this.initHTML();
		defineShopProductWidget();
		this.shadowRoot!.replaceChildren();
		for (const shopProduct of products) {
			createElement('shop-product-widget', {
				parent: this.shadowRoot,
				elementCreated: (element: HTMLElement) => (element as HTMLShopProductWidgetElement).setProduct(shopProduct),
			});
		}
	}
}
