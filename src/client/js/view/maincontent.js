import { createElement, hide, show } from 'harmony-ui';

import mainContentCSS from '../../css/maincontent.css';
import { ProductsPage } from './productspage';
import { PAGE_TYPE_CART, PAGE_TYPE_CHECKOUT, PAGE_TYPE_LOGIN, PAGE_TYPE_ORDER, PAGE_TYPE_OTHER, PAGE_TYPE_PRODUCTS, PAGE_TYPE_SHOP, PAGE_TYPE_UNKNOWN } from '../constants';

export class MainContent {
	#htmlElement;
	#productsPage = new ProductsPage();

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: mainContentCSS,
			childs: [
				this.#productsPage.htmlElement,
			],
		});
		this.setActivePage(PAGE_TYPE_UNKNOWN);
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setActivePage(pageType) {
		hide(this.#productsPage.htmlElement);

		switch (pageType) {
			case PAGE_TYPE_UNKNOWN:
				break;
			case PAGE_TYPE_PRODUCTS:
				show(this.#productsPage.htmlElement);
				break;
			case PAGE_TYPE_SHOP:
				throw 'TODO: PAGE_TYPE_SHOP';
				break;
			case PAGE_TYPE_CART:
				throw 'TODO: PAGE_TYPE_CART';
				break;
			case PAGE_TYPE_CHECKOUT:
				throw 'TODO: PAGE_TYPE_CHECKOUT';
				break;
			case PAGE_TYPE_LOGIN:
				throw 'TODO: PAGE_TYPE_LOGIN';
				break;
			case PAGE_TYPE_OTHER:
				throw 'TODO: PAGE_TYPE_OTHER';
				break;
			case PAGE_TYPE_ORDER:
				throw 'TODO: PAGE_TYPE_ORDER';
				break;
			default:
				throw `Unknown page type ${pageType}`;
		}
	}

	setProducts(products) {
		this.#productsPage.setProducts(products);
	}
}
