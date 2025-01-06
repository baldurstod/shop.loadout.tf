import { createShadowRoot } from 'harmony-ui';
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
import { OrderPage } from './orderpage';
import { ShopElement } from './shopelement';

export class MainContent extends ShopElement {
	#cartPage = new CartPage();
	#checkoutPage = new CheckoutPage();
	#contactPage = new ContactPage();
	#cookiesPage = new CookiesPage();
	#favoritesPage = new FavoritesPage();
	#privacyPage = new PrivacyPage();
	#productPage = new ProductPage();
	#productsPage = new ProductsPage();
	#orderPage = new OrderPage();

	initHTML() {
		if (this.shadowRoot) {
			return;
		}

		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: mainContentCSS,
		});
		this.setActivePage(PAGE_TYPE_UNKNOWN);
	}

	setActivePage(pageType: PageType, pageSubType?: PageSubType) {
		this.initHTML();
		this.shadowRoot?.replaceChildren();

		switch (pageType) {
			case PAGE_TYPE_UNKNOWN:
				break;
			case PAGE_TYPE_CART:
				this.shadowRoot?.append(this.#cartPage.getHTML());
				break;
			case PAGE_TYPE_CHECKOUT:
				this.#checkoutPage.setCheckoutStage(pageSubType ?? PageSubType.CheckoutInit);
				this.shadowRoot?.append(this.#checkoutPage.getHTML());
				break;
			case PAGE_TYPE_LOGIN:
				throw 'TODO: PAGE_TYPE_LOGIN';
				break;
			case PAGE_TYPE_ORDER:
				this.shadowRoot?.append(this.#orderPage.getHTML());
				break;
			case PAGE_TYPE_PRODUCTS:
				this.shadowRoot?.append(this.#productsPage.getHTML());
				break;
			case PAGE_TYPE_COOKIES:
				this.shadowRoot?.append(this.#cookiesPage.getHTML());
				break;
			case PAGE_TYPE_PRIVACY:
				this.shadowRoot?.append(this.#privacyPage.getHTML());
				break;
			case PAGE_TYPE_CONTACT:
				this.shadowRoot?.append(this.#contactPage.getHTML());
				break;
			case PAGE_TYPE_PRODUCT:
				this.shadowRoot?.append(this.#productPage.getHTML());
				break;
			case PAGE_TYPE_FAVORITES:
				this.shadowRoot?.append(this.#favoritesPage.getHTML());
				break;
			default:
				throw `Unknown page type ${pageType}`;
		}
	}

	setProduct(product: Product) {
		this.#productPage.setProduct(product);
	}

	setCheckoutOrder(order: Order) {
		this.#checkoutPage.setOrder(order);
	}

	setOrder(order: Order) {
		this.#orderPage.setOrder(order);
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
