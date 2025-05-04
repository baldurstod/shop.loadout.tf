import { createShadowRoot } from 'harmony-ui';
import { CartPage } from './cartpage';
import { CheckoutPage } from './checkoutpage';
import { ContactPage } from './contactpage';
import { CookiesPage } from './cookiespage';
import { FavoritesPage } from './favoritespage';
import { PrivacyPage } from './privacypage';
import { ProductPage } from './productpage';
import { ProductsPage } from './productspage';
import { PageSubType, PageType } from '../constants';
import mainContentCSS from '../../css/maincontent.css';
import { Product } from '../model/product';
import { Order } from '../model/order';
import { Cart } from '../model/cart';
import { Countries } from '../model/countries';
import { OrderPage } from './orderpage';
import { ShopElement } from './shopelement';
import { LogoutPage } from './logoutpage';
import { LoginPage } from './loginpage';

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
	#loginPage = new LoginPage();
	#logoutPage = new LogoutPage();

	initHTML() {
		if (this.shadowRoot) {
			return;
		}

		this.shadowRoot = createShadowRoot('section', {
			adoptStyle: mainContentCSS,
		});
		this.setActivePage(PageType.Unknown);
	}

	setActivePage(pageType: PageType, pageSubType?: PageSubType) {
		this.initHTML();
		this.shadowRoot?.replaceChildren();

		switch (pageType) {
			case PageType.Unknown:
				break;
			case PageType.Cart:
				this.shadowRoot?.append(this.#cartPage.getHTML());
				break;
			case PageType.Checkout:
				this.#checkoutPage.setCheckoutStage(pageSubType ?? PageSubType.CheckoutInit);
				this.shadowRoot?.append(this.#checkoutPage.getHTML());
				break;
			case PageType.Login:
				this.shadowRoot?.append(this.#loginPage.getHTML());
				break;
			case PageType.Logout:
				this.shadowRoot?.append(this.#logoutPage.getHTML());
				break;
			case PageType.Order:
				this.shadowRoot?.append(this.#orderPage.getHTML());
				break;
			case PageType.Products:
				this.shadowRoot?.append(this.#productsPage.getHTML());
				break;
			case PageType.Cookies:
				this.shadowRoot?.append(this.#cookiesPage.getHTML());
				break;
			case PageType.Privacy:
				this.shadowRoot?.append(this.#privacyPage.getHTML());
				break;
			case PageType.Contact:
				this.shadowRoot?.append(this.#contactPage.getHTML());
				break;
			case PageType.Product:
				this.shadowRoot?.append(this.#productPage.getHTML());
				break;
			case PageType.Favorites:
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
		this.#productPage.refreshFavorite();
	}

	setCart(cart: Cart) {
		this.#cartPage.setCart(cart);
	}

	setCountries(countries: Countries) {
		this.#checkoutPage.setCountries(countries);
	}
}
