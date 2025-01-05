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
import { Cart } from '../model/cart';
import { Countries } from '../model/countries';

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
				this.#checkoutPage.getHTML(),
				this.#contactPage.getHTML(),
				this.#cookiesPage.getHTML(),
				this.#favoritesPage.getHTML(),
				this.#privacyPage.getHTML(),
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
		hide(this.#checkoutPage.getHTML());
		hide(this.#contactPage.getHTML());
		hide(this.#cookiesPage.getHTML());
		hide(this.#favoritesPage.getHTML());
		hide(this.#privacyPage.getHTML());
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
				show(this.#checkoutPage.getHTML());
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
				show(this.#cookiesPage.getHTML());
				break;
			case PAGE_TYPE_PRIVACY:
				show(this.#privacyPage.getHTML());
				break;
			case PAGE_TYPE_CONTACT:
				show(this.#contactPage.getHTML());
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

	setCart(cart: Cart) {
		this.#cartPage.setCart(cart);
	}

	setCountries(countries: Countries) {
		this.#checkoutPage.setCountries(countries);
	}
}
