import { createElement } from 'harmony-ui';

import favoritesPageCSS from '../../css/favoritespage.css';
import { defineShopProductWidget } from './components/shopproductwidget';

export class FavoritesPage {
	#htmlElement;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: favoritesPageCSS,
		});
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setFavorites(favorites) {
		this.#htmlElement.innerHTML = '';
		defineShopProductWidget();
		for (const shopProduct of favorites) {
			createElement('shop-product-widget', {
				parent: this.#htmlElement,
				elementCreated: element => element.setProduct(shopProduct),
			});
		}
	}
}
