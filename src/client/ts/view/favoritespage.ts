import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import favoritesPageCSS from '../../css/favoritespage.css';
import { defineShopProductWidget, HTMLShopProductWidgetElement } from './components/shopproductwidget';
import { Product } from '../model/product';

export class FavoritesPage {
	#shadowRoot?: ShadowRoot;

	#initHTML() {
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyle: favoritesPageCSS,
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}

	setFavorites(favorites: Array<Product>) {
		if (!this.#shadowRoot) {
			this.#initHTML();
		}
		this.#shadowRoot!.innerHTML = '';
		defineShopProductWidget();
		for (const shopProduct of favorites) {
			createElement('shop-product-widget', {
				parent: this.#shadowRoot,
				elementCreated: (element: HTMLShopProductWidgetElement) => element.setProduct(shopProduct),
			});
		}
	}
}
