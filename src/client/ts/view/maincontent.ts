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
import { Order } from '../model/order';

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
				this.#cartPage.getHTML(),
				this.#checkoutPage.htmlElement,
				this.#contactPage.htmlElement,
				this.#cookiesPage.htmlElement,
				this.#favoritesPage.getHTML(),
				this.#privacyPage.htmlElement,
				this.#productPage.getHTML(),
				this.#productsPage.getHTML(),
			],
		});
		this.setActivePage(PAGE_TYPE_UNKNOWN);
		return this.#shadowRoot.host;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}

	setActivePage(pageType: PageType, pageSubType?: PageSubType) {
		hide(this.#cartPage.getHTML());
		hide(this.#checkoutPage.htmlElement);
		hide(this.#contactPage.htmlElement);
		hide(this.#cookiesPage.htmlElement);
		hide(this.#favoritesPage.getHTML());
		hide(this.#privacyPage.htmlElement);
		hide(this.#productPage.getHTML());
		hide(this.#productsPage.getHTML());

		switch (pageType) {
			case PAGE_TYPE_UNKNOWN:
				break;
			case PAGE_TYPE_CART:
				show(this.#cartPage.getHTML());
				break;
			case PAGE_TYPE_CHECKOUT:
				this.#checkoutPage.setCheckoutStage(pageSubType ?? PageSubType.CheckoutInit);
				show(this.#checkoutPage.htmlElement);
				break;
			case PAGE_TYPE_LOGIN:
				throw 'TODO: PAGE_TYPE_LOGIN';
				break;
			case PAGE_TYPE_ORDER:
				throw 'TODO: PAGE_TYPE_ORDER';
				break;
			case PAGE_TYPE_PRODUCTS:
				show(this.#productsPage.getHTML());
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
				show(this.#favoritesPage.getHTML());
				break;
			default:
				throw `Unknown page type ${pageType}`;
		}
	}

	setProduct(product: Product) {
		this.#productPage.setProduct(product);
	}

	setOrder(order: Order) {
		this.#checkoutPage.setOrder(order);
	}

	setProducts(products: Array<Product>) {
		this.#productsPage.setProducts(products);
	}

	setFavorites(favorites: Array<Product>) {
		this.#favoritesPage.setFavorites(favorites);
	}

	setCart(cart) {
		this.#cartPage.setCart(cart);
	}

	setCountries(countries) {
		this.#checkoutPage.setCountries(countries);
	}
}
