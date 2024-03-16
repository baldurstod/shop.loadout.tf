import { createElement, hide, show } from 'harmony-ui';
import { ContactPage } from './contactpage.js';
import { PrivacyPage } from './privacypage.js';
import { ProductsPage } from './productspage.js';
import { PAGE_TYPE_CART, PAGE_TYPE_CHECKOUT, PAGE_TYPE_CONTACT, PAGE_TYPE_COOKIES, PAGE_TYPE_LOGIN, PAGE_TYPE_ORDER, PAGE_TYPE_PRIVACY, PAGE_TYPE_PRODUCTS, PAGE_TYPE_SHOP, PAGE_TYPE_UNKNOWN } from '../constants.js';

import mainContentCSS from '../../css/maincontent.css';

export class MainContent {
	#htmlElement;
	#contactPage = new ContactPage();
	#privacyPage = new PrivacyPage();
	#productsPage = new ProductsPage();

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: mainContentCSS,
			childs: [
				this.#contactPage.htmlElement,
				this.#privacyPage.htmlElement,
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
		hide(this.#contactPage.htmlElement);
		hide(this.#privacyPage.htmlElement);
		hide(this.#productsPage.htmlElement);

		switch (pageType) {
			case PAGE_TYPE_UNKNOWN:
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
			case PAGE_TYPE_ORDER:
				throw 'TODO: PAGE_TYPE_ORDER';
				break;
			case PAGE_TYPE_PRODUCTS:
				show(this.#productsPage.htmlElement);
				break;
			case PAGE_TYPE_COOKIES:
				throw 'TODO: PAGE_TYPE_COOKIES';
				break;
			case PAGE_TYPE_PRIVACY:
				show(this.#privacyPage.htmlElement);
				break;
			case PAGE_TYPE_CONTACT:
				show(this.#contactPage.htmlElement);
				break;
			default:
				throw `Unknown page type ${pageType}`;
		}
	}

	setProducts(products) {
		this.#productsPage.setProducts(products);
	}
}
