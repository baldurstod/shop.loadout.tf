import { createElement } from 'harmony-ui';
export { ShopProductWidgetElement } from './components/shopproductwidget.js';

import favoritesPageCSS from '../../css/favoritespage.css';

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
		for (const shopProduct of favorites) {
			createElement('shop-product-widget', {
				parent: this.#htmlElement,
				elementCreated: element => element.setProduct(shopProduct),
			});
		}
	}
}
