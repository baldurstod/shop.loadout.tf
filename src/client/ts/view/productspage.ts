import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import productsPageCSS from '../../css/productspage.css';
import { Product } from '../model/product';
import { defineShopProductWidget, HTMLShopProductWidgetElement } from './components/shopproductwidget';
import { ShopElement } from './shopelement';

export class ProductsPage extends ShopElement {

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: productsPageCSS,
		});
		I18n.observeElement(this.shadowRoot);
	}

	setProducts(products: Product[] = []): void {
		this.initHTML();
		defineShopProductWidget();
		this.shadowRoot!.replaceChildren();
		for (const shopProduct of products) {
			createElement('shop-product-widget', {
				parent: this.shadowRoot,
				elementCreated: (element: Element) => (element as HTMLShopProductWidgetElement).setProduct(shopProduct),
			});
		}
	}
}
