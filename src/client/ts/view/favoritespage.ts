import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import favoritesPageCSS from '../../css/favoritespage.css';
import { defineShopProductWidget, HTMLShopProductWidgetElement } from './components/shopproductwidget';
import { Product } from '../model/product';
import { ShopElement } from './shopelement';

export class FavoritesPage extends ShopElement {
	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: favoritesPageCSS,
		});
		I18n.observeElement(this.shadowRoot);
	}

	setFavorites(favorites: Array<Product>) {
		this.initHTML();
		this.shadowRoot!.replaceChildren();
		defineShopProductWidget();
		for (const shopProduct of favorites) {
			createElement('shop-product-widget', {
				parent: this.shadowRoot,
				elementCreated: (element: HTMLElement) => (element as HTMLShopProductWidgetElement).setProduct(shopProduct),
			});
		}
	}
}
