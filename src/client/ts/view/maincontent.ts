import { createElement, hide, show } from 'harmony-ui';
import { CartPage } from './cartpage.js';
import { CheckoutPage } from './checkoutpage.js';
import { ContactPage } from './contactpage.js';
import { CookiesPage } from './cookiespage.js';
import { FavoritesPage } from './favoritespage.js';
import { PrivacyPage } from './privacypage.js';
import { ProductPage } from './productpage.js';
import { ProductsPage } from './productspage.js';
import { PAGE_TYPE_CART, PAGE_TYPE_CHECKOUT, PAGE_TYPE_CONTACT, PAGE_TYPE_COOKIES, PAGE_TYPE_FAVORITES, PAGE_TYPE_LOGIN, PAGE_TYPE_ORDER, PAGE_TYPE_PRIVACY, PAGE_TYPE_PRODUCT, PAGE_TYPE_PRODUCTS, PAGE_TYPE_UNKNOWN } from '../constants.js';

import mainContentCSS from '../../css/maincontent.css';

export class MainContent {
	#htmlElement;
	#cartPage = new CartPage();
	#checkoutPage = new CheckoutPage();
	#contactPage = new ContactPage();
	#cookiesPage = new CookiesPage();
	#favoritesPage = new FavoritesPage();
	#privacyPage = new PrivacyPage();
	#productPage = new ProductPage();
	#productsPage = new ProductsPage();

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: mainContentCSS,
			childs: [
				this.#cartPage.htmlElement,
				this.#checkoutPage.htmlElement,
				this.#contactPage.htmlElement,
				this.#cookiesPage.htmlElement,
				this.#favoritesPage.htmlElement,
				this.#privacyPage.htmlElement,
				this.#productPage.htmlElement,
				this.#productsPage.htmlElement,
			],
		});
		this.setActivePage(PAGE_TYPE_UNKNOWN);
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	setActivePage(pageType, pageSubType?) {
		hide(this.#cartPage.htmlElement);
		hide(this.#checkoutPage.htmlElement);
		hide(this.#contactPage.htmlElement);
		hide(this.#cookiesPage.htmlElement);
		hide(this.#favoritesPage.htmlElement);
		hide(this.#privacyPage.htmlElement);
		hide(this.#productPage.htmlElement);
		hide(this.#productsPage.htmlElement);

		switch (pageType) {
			case PAGE_TYPE_UNKNOWN:
				break;
			case PAGE_TYPE_CART:
				show(this.#cartPage.htmlElement);
				break;
			case PAGE_TYPE_CHECKOUT:
				this.#checkoutPage.setCheckoutStage(pageSubType);
				show(this.#checkoutPage.htmlElement);
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
				show(this.#cookiesPage.htmlElement);
				break;
			case PAGE_TYPE_PRIVACY:
				show(this.#privacyPage.htmlElement);
				break;
			case PAGE_TYPE_CONTACT:
				show(this.#contactPage.htmlElement);
				break;
			case PAGE_TYPE_PRODUCT:
				show(this.#productPage.htmlElement);
				break;
			case PAGE_TYPE_FAVORITES:
				show(this.#favoritesPage.htmlElement);
				break;
			default:
				throw `Unknown page type ${pageType}`;
		}
	}

	setProduct(product) {
		this.#productPage.setProduct(product);
	}

	setOrder(order) {
		this.#checkoutPage.setOrder(order);
	}

	setProducts(products) {
		this.#productsPage.setProducts(products);
	}

	setFavorites(favorites) {
		this.#favoritesPage.setFavorites(favorites);
	}

	setCart(cart) {
		this.#cartPage.setCart(cart);
	}

	setCountries(countries) {
		this.#checkoutPage.setCountries(countries);
	}
}
