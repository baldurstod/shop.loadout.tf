import { createElement, createShadowRoot, hide, show } from 'harmony-ui';
import { CartPage } from './cartpage';
import { CheckoutPage } from './checkoutpage';
import { ContactPage } from './contactpage';
import { CookiesPage } from './cookiespage';
import { FavoritesPage } from './favoritespage';
import { PrivacyPage } from './privacypage';
import { ProductPage } from './productpage';
import { ProductsPage } from './productspage';
import { PAGE_TYPE_CART, PAGE_TYPE_CHECKOUT, PAGE_TYPE_CONTACT, PAGE_TYPE_COOKIES, PAGE_TYPE_FAVORITES, PAGE_TYPE_LOGIN, PAGE_TYPE_ORDER, PAGE_TYPE_PRIVACY, PAGE_TYPE_PRODUCT, PAGE_TYPE_PRODUCTS, PAGE_TYPE_UNKNOWN, PageSubType, PageType } from '../constants';
import mainContentCSS from '../../css/maincontent.css';
import { Product } from '../model/product';

export class MainContent {
	#shadowRoot?: ShadowRoot;
	#cartPage = new CartPage();
	#checkoutPage = new CheckoutPage();
	#contactPage = new ContactPage();
	#cookiesPage = new CookiesPage();
	#favoritesPage = new FavoritesPage();
	#privacyPage = new PrivacyPage();
	#productPage = new ProductPage();
	#productsPage = new ProductsPage();

	#initHTML() {
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyle: mainContentCSS,
			childs: [
				this.#cartPage.htmlElement,
				this.#checkoutPage.htmlElement,
				this.#contactPage.htmlElement,
				this.#cookiesPage.htmlElement,
				this.#favoritesPage.htmlElement,
				this.#privacyPage.htmlElement,
				this.#productPage.getHTML(),
				this.#productsPage.htmlElement,
			],
		});
		this.setActivePage(PAGE_TYPE_UNKNOWN);
		return this.#shadowRoot.host;
	}

	get htmlElement() {
		throw 'use getHTML';
	}

	getHTML() {
		return this.#shadowRoot?.host ?? this.#initHTML();
	}

	setActivePage(pageType: PageType, pageSubType?: PageSubType) {
		hide(this.#cartPage.htmlElement);
		hide(this.#checkoutPage.htmlElement);
		hide(this.#contactPage.htmlElement);
		hide(this.#cookiesPage.htmlElement);
		hide(this.#favoritesPage.htmlElement);
		hide(this.#privacyPage.htmlElement);
		hide(this.#productPage.getHTML());
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
				show(this.#productPage.getHTML());
				break;
			case PAGE_TYPE_FAVORITES:
				show(this.#favoritesPage.htmlElement);
				break;
			default:
				throw `Unknown page type ${pageType}`;
		}
	}

	setProduct(product: Product) {
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
