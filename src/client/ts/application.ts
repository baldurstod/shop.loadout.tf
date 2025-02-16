import { addNotification } from 'harmony-browser-utils';
import { themeCSS } from 'harmony-css';
import { createElement, I18n, documentStyle, defineHarmonyCopy, defineHarmonySwitch, defineHarmonyPalette, defineHarmonySlideshow, createShadowRoot } from 'harmony-ui';
import { getShopProduct } from './shopproducts';
import { PAYPAL_APP_CLIENT_ID, BROADCAST_CHANNEL_NAME, PAGE_TYPE_CART, PAGE_TYPE_CHECKOUT, PAGE_TYPE_PRODUCTS, PAGE_TYPE_COOKIES, PAGE_TYPE_PRIVACY, PAGE_TYPE_CONTACT, PAGE_TYPE_LOGIN, PAGE_TYPE_ORDER, PAGE_TYPE_PRODUCT, PAGE_TYPE_FAVORITES, PAGE_SUBTYPE_CHECKOUT_INIT, PAGE_SUBTYPE_CHECKOUT_ADDRESS, PAGE_SUBTYPE_CHECKOUT_SHIPPING, PAGE_SUBTYPE_CHECKOUT_PAYMENT, PAGE_SUBTYPE_CHECKOUT_COMPLETE, PAGE_SUBTYPE_SHOP_PRODUCT, PageType, PageSubType } from './constants';
import { Controller } from './controller';
import { Footer } from './view/footer';
import { MainContent } from './view/maincontent';
import { Toolbar } from './view/toolbar';
//import { OrderSummary } from './view/ordersummary';
import { Cart } from './model/cart';
import { Order } from './model/order';
import { Product } from './model/product';
import '../css/item.css';
import '../css/order.css';
import '../css/shop.css';
import '../css/vars.css';
import applicationCSS from '../css/application.css';
import htmlCSS from '../css/html.css';
import english from '../json/i18n/english.json';
import { fetchApi } from './fetchapi';
import { ServerAPI } from './serverapi';
import { EVENT_CART_COUNT, EVENT_DECREASE_FONT_SIZE, EVENT_FAVORITES_COUNT, EVENT_INCREASE_FONT_SIZE, EVENT_NAVIGATE_TO, EVENT_REFRESH_CART, EVENT_SEND_CONTACT, EVENT_SEND_CONTACT_ERROR } from './controllerevents';
import { Countries } from './model/countries';
import { BroadcastMessage } from './enums';
import { defineShopProduct, HTMLShopProductElement } from './view/components/shopproduct';
import { JSONObject } from './types';
import { FavoritesResponse } from './responses/favorites';
import { CountriesResponse } from './responses/countries';
import { InitCheckoutResponse, OrderJSON, OrderResponse, SetShippingAddressResponse, SetShippingMethodResponse } from './responses/order';
import { GetProductsResponse } from './responses/products';
import { AddProductResponse, GetCartResponse } from './responses/cart';
import { favoritesCount, getFavorites, setFavorites, toggleFavorite } from './favorites';

const REFRESH_PRODUCT_PAGE_DELAY = 20000;

documentStyle(htmlCSS);
documentStyle(themeCSS);

const TESTING = false;



class Application {
	#appToolbar = new Toolbar();
	#appContent = new MainContent();
	#appFooter = new Footer();
	#pageType: PageType = PageType.Unknown;
	#pageSubType: PageSubType = PageSubType.Unknown;
	#order: Order | null = null;
	//#orderSummary = new OrderSummary();
	//#favorites: Set<string> = new Set();
	#broadcastChannel = new BroadcastChannel(BROADCAST_CHANNEL_NAME);
	#paymentCompleteDetails: { order: Order }/*TODO: improve type*/ | null = null;
	//#htmlColumnCart;
	//#htmlColumnCartVisible = false;
	#htmlShopProduct?: HTMLShopProductElement;
	#orderId?: string;
	#fontSize = 16;
	#refreshPageTimeout?: ReturnType<typeof setTimeout>;
	#cart = new Cart();
	#countries = new Countries();

	constructor() {
		I18n.setOptions({ translations: [english] });
		I18n.start();

		Controller.addEventListener('addtocart', (event: Event) => this.#addToCart((event as CustomEvent).detail.product, (event as CustomEvent).detail.quantity));
		Controller.addEventListener('setquantity', (event: Event) => this.#setQuantity((event as CustomEvent).detail.id, (event as CustomEvent).detail.quantity));
		Controller.addEventListener(EVENT_NAVIGATE_TO, (event: Event) => this.#navigateTo((event as CustomEvent).detail.url, (event as CustomEvent).detail.replaceSate));
		Controller.addEventListener('pushstate', (event: Event) => this.#pushState((event as CustomEvent).detail.url));
		Controller.addEventListener('replacestate', (event: Event) => this.#replaceState((event as CustomEvent).detail.url));
		Controller.addEventListener('addnotification', (event: Event) => addNotification((event as CustomEvent).detail.content, (event as CustomEvent).detail.type));
		Controller.addEventListener('paymentcomplete', (event: Event) => this.#onPaymentComplete((event as CustomEvent).detail));
		Controller.addEventListener('favorite', (event: Event) => this.#favorite((event as CustomEvent).detail.productId));
		Controller.addEventListener('schedulerefreshproductpage', () => this.#scheduleRefreshProductPage());
		Controller.addEventListener(EVENT_REFRESH_CART, () => this.#refreshCart());
		this.#initListeners();

		new BroadcastChannel(BROADCAST_CHANNEL_NAME).addEventListener('message', (event) => {
			this.#processMessage(event);
		});

		//this.theme = 'light';

		this.#initPage();
		this.#initSession();
		this.#startup();
		this.#initFavorites();
		this.#initCountries();
		addEventListener('popstate', event => this.#startup(event.state ?? {}));
		this.#loadCart();
	}

	#initListeners() {
		Controller.addEventListener(EVENT_INCREASE_FONT_SIZE, () => this.#changeFontSize(1));
		Controller.addEventListener(EVENT_DECREASE_FONT_SIZE, () => this.#changeFontSize(-1));
		Controller.addEventListener(EVENT_SEND_CONTACT, (event: Event) => this.#sendContact((event as CustomEvent).detail));
	}

	#changeFontSize(change: number) {
		let size = this.#fontSize;
		if (change > 0) {
			size *= 1.1;
		} else {
			size /= 1.1;
		}
		this.#broadcastChannel.postMessage({ action: BroadcastMessage.FontSizeChanged, fontSize: size });
	}

	#setFontSize(size: number) {
		this.#fontSize = size;
		document.documentElement.style.setProperty('--font-size', `${this.#fontSize}px`);
	}


	async #startup(historyState = {}) {
		this.#restoreHistoryState(historyState);
		let pathname = document.location.pathname;
		this.#pageSubType = PageSubType.Unknown;
		switch (true) {
			case pathname.includes('@cart'):
				this.#pageType = PAGE_TYPE_CART;
				break;
			case pathname.includes('@products'):
				this.#pageType = PAGE_TYPE_PRODUCTS;
				this.#displayProducts();
				break;
			case pathname.includes('@favorites'):
				this.#pageType = PAGE_TYPE_FAVORITES;
				break;
			case pathname.includes('@product'):
				this.#pageType = PAGE_TYPE_PRODUCT;
				await this.#initProductFromUrl();
				break;
			case pathname.includes('@checkout'):
				this.#pageType = PAGE_TYPE_CHECKOUT;
				switch (document.location.hash) {
					case '':
						this.#pageSubType = PAGE_SUBTYPE_CHECKOUT_INIT;
						this.#initCheckout();
						break;
					case '#address':
						this.#pageSubType = PAGE_SUBTYPE_CHECKOUT_ADDRESS;
						this.#initAddress();
						break;
					case '#shipping':
						this.#pageSubType = PAGE_SUBTYPE_CHECKOUT_SHIPPING;
						this.#initShipping();
						break;
					case '#payment':
						this.#pageSubType = PAGE_SUBTYPE_CHECKOUT_PAYMENT;
						this.#initPayment();
						break;
					case '#complete':
						this.#pageSubType = PAGE_SUBTYPE_CHECKOUT_COMPLETE;
						this.#paymentComplete();
						break;
				}
				break;
			case pathname.includes('@login'):
				this.#pageType = PAGE_TYPE_LOGIN;
				this.#viewLoginPage();
				break;
			case pathname.includes('@cookies'):
				this.#pageType = PAGE_TYPE_COOKIES;
				break;
			case pathname.includes('@privacy'):
				this.#pageType = PAGE_TYPE_PRIVACY;
				break;
			case pathname.includes('@contact'):
				this.#pageType = PAGE_TYPE_CONTACT;
				break;
			case pathname.includes('@order'):
				this.#pageType = PAGE_TYPE_ORDER;
				this.#initOrderFromUrl();
				break;
			default:
				this.#navigateTo('/@products');
				break;
		}

		this.#appContent.setActivePage(this.#pageType, this.#pageSubType);
	}

	async #initFavorites() {
		const { requestId, response } = await fetchApi({
			action: 'get-favorites',
			version: 1,
		}) as { requestId: string, response: FavoritesResponse };
		if (response?.success) {
			setFavorites(response.result?.favorites);

			this.#countFavorites();
			this.#broadcastChannel.postMessage({ action: BroadcastMessage.FavoritesChanged, favorites: getFavorites() });
		}
	}

	async #initCountries() {
		const { requestId, response } = await fetchApi({
			action: 'get-countries',
			version: 1,
		}) as { requestId: string, response: CountriesResponse };
		if (response?.success && response.result?.countries) {
			this.#countries.fromJSON(response.result.countries);
			this.#appContent.setCountries(this.#countries);
		}
	}

	async #favorite(productId: string) {
		await fetchApi({
			action: 'set-favorite',
			version: 1,
			params: {
				product_id: productId,
				is_favorite: toggleFavorite(productId),
			},
		});

		this.#broadcastChannel.postMessage({ action: BroadcastMessage.FavoritesChanged, favorites: getFavorites() });
		this.#countFavorites();
	}

	#countFavorites() {
		Controller.dispatchEvent(new CustomEvent(EVENT_FAVORITES_COUNT, { detail: favoritesCount() }));
	}

	async #addToCart(productId: string, quantity = 1) {
		if (TESTING) {
			console.log(productId, quantity);
		}

		const { requestId, response } = await fetchApi({
			action: 'add-product',
			version: 1,
			params: {
				product_id: productId,
				quantity: quantity,
			},
		}) as { requestId: string, response: AddProductResponse };

		if (response.success && response.result?.cart) {
			this.#cart.fromJSON(response.result.cart);
			this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartChanged, cart: this.#cart.toJSON() });
		}

	}

	async #setQuantity(productId: string, quantity = 1) {
		if (TESTING) {
			console.log(productId, quantity);
		}

		const { requestId, response } = await fetchApi({
			action: 'set-product-quantity',
			version: 1,
			params: {
				product_id: productId,
				quantity: quantity,
			},
		}) as { requestId: string, response: AddProductResponse };

		if (response.success && response.result?.cart) {
			this.#cart.fromJSON(response.result.cart);
			this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartChanged, cart: this.#cart.toJSON() });
		}
	}

	async #refreshCart() {
		this.#appContent.setCart(this.#cart);
	}

	async #refreshFavorites() {
		const favorites: Array<Product> = [];

		for (const productID of getFavorites()) {
			const product = await getShopProduct(productID);
			if (product) {
				favorites.push(product);
			}
		}

		this.#appContent.setFavorites(favorites);
	}

	async #initProductFromUrl() {
		let result = /@product\/([^\/]*)/i.exec(document.location.pathname);
		if (result) {
			await this.#initProductPage(result[1]);
		}

		this.#refreshCart();
	}

	async #initOrderFromUrl() {
		let result = /@order\/([^\/]*)/i.exec(document.location.pathname);
		if (result) {
			this.#loadCart();
			await this.#initOrderPage(result[1]);
		}
	}

	async #viewLoginPage() {
		/*
		let htmlLoginPage = createElement('div', { class: 'shop-login-page' });
		htmlLoginPage.append(createElement('div', { id: 'paypal-login-container' }));
		await loadScript('https://www.paypalobjects.com/js/external/api.js');

		paypal.use(['login'], function (login) {
			login.render({
				"appid": PAYPAL_APP_CLIENT_ID,
				"authend": "sandbox",
				"scopes": "profile email address",
				"containerid": "paypal-login-container",
				"responseType": "code",
				"locale": "en-us",
				"buttonType": "LWP",
				"buttonShape": "pill",
				"buttonSize": "lg",
				"fullPage": "true",
				"returnurl": document.location.origin + '/paypalredirect',
				"nonce": "12345678"
			});
		});

		this.htmlContent.append(htmlLoginPage);
		*/
	}

	async #sendContact(detail: { subject: string, email: string, content: string, }) {
		const { requestId, response } = await fetchApi({
			action: 'send-contact',
			version: 1,
			params: {
				subject: detail.subject,
				email: detail.email,
				content: detail.content,
			},
		});


		if (response?.success) {
			Controller.dispatchEvent(new CustomEvent('addnotification', { detail: { type: 'success', content: createElement('span', { i18n: '#message_successfully_sent' }) } }));
			//detail.callback(true);
		} else {
			Controller.dispatchEvent(new CustomEvent('addnotification', { detail: { type: 'error', content: createElement('span', { i18n: '#error_while_sending_message' }) } }));
			Controller.dispatchEvent(new CustomEvent(EVENT_SEND_CONTACT_ERROR));
		}
	}

	#displayCheckout() {
		return;
		/*
		let htmlCheckoutPage = createElement('div', { class: 'shop-checkout-page' });

		htmlCheckoutPage.append(this.#orderSummary.html);
		this.htmlContent.append(htmlCheckoutPage);
		this.#refreshOrder();
		*/
	}

	#displayPaymentComplete() {
		/*
		let cartItems = [];

		let order = this.#paymentCompleteDetails.order;
		let paymentDetail = this.#paymentCompleteDetails.paymentDetail;
		if (TESTING) {
			console.log(order, paymentDetail);
		}

		let orderSummary = new OrderSummary();
		orderSummary.summary = order;

		let htmlPaymentCompletePage = createElement('div', {
			class: 'shop-payment-complete-page',
			childs: [
				createElement('div', {
					class: 'shop-payment-complete-header',
					i18n: '#payment_complete',
				}),
				/*createElement('div', {
					childs: [
						createElement('label', {
							i18n: '#order_id',
						}),
						createElement('span', {
							innerHTML: order.id,
						}),
					]
				}),* /
				createElement('label-property', {
					label: '#order_id',
					property: order.id,
				}),
				transactionSummary(paymentDetail),
				createElement('div', {
					class: 'shop-payment-complete-content',
					childs: [
						createElement('div', {
							class: 'shop-payment-complete-cart',
							child: orderSummary.html,
						}),

					]
				}),
			]
		});


		this.htmlContent.append(htmlPaymentCompletePage);
		*/
	}

	async #initCheckout() {
		const { requestId, response } = await fetchApi({
			action: 'init-checkout',
			version: 1
		}) as { requestId: string, response: InitCheckoutResponse };
		if (response?.success && response.result?.order) {
			const order = new Order();
			order.fromJSON(response.result.order);
			this.#appContent.setCheckoutOrder(order);

			this.#order = order;

			//this.#orderSummary.setOrder(order);
			this.#orderId = order.id;

			this.#navigateTo('/@checkout#address', true);

		} else {
			Controller.dispatchEvent(new CustomEvent('addnotification', { detail: { type: 'error', content: createElement('span', { i18n: '#failed_to_init_order' }) } }));
			return false;
		}
	}

	async #refreshOrder() {
		const response = await fetch('/api', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				action: 'getorder',
				version: 1,
				orderId: this.#orderId,
			}),
		});

		let json = await response.json();
		if (json?.success) {
			let order = new Order();
			order.fromJSON(json.result);
			//this.#orderSummary.setOrder(order);
		}
	}

	#initAddress() {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}

		this.#displayCheckout();
	}

	async #initShipping() {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}

		await this.#sendShippingAddress();

		/*let response = await fetch('/gettaxrates');
		let json = await response.json();
		if (json?.success) {
			//this.#navigateTo('/@checkout#shipping');
		}*/
		this.#displayCheckout();
	}

	async #sendShippingAddress() {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}

		const { requestId, response } = await fetchApi({
			action: 'set-shipping-address',
			version: 1,
			params: {
				shipping_address: this.#order.shippingAddress,
				//orderId: this.#orderId
			},
		}) as { requestId: string, response: SetShippingAddressResponse };


		/*
		const { requestId, response } = await fetchApi({
			action: 'set-product-quantity',
			version: 1,
			params: {
				product_id: productId,
				quantity: quantity,
				//cart: this.#cart.toJSON(),
			},
		});
		*/



		if (response?.success && response.result?.order) {
			this.#order.fromJSON(response.result.order);
			this.#appContent.setCheckoutOrder(this.#order);
			return true;
		}
	}

	async #sendShippingMethod() {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}

		const { requestId, response } = await fetchApi({
			action: 'set-shipping-method',
			version: 1,
			params: {
				method: this.#order.shippingMethod,
			},
		}) as { requestId: string, response: SetShippingMethodResponse };

		if (response?.success && response.result?.order) {
			this.#order.fromJSON(response.result.order);
			this.#appContent.setCheckoutOrder(this.#order);
			return true;
		}
	}

	async #initPayment() {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}

		this.#paymentCompleteDetails = null;

		const shippingOK = await this.#sendShippingMethod();
		if (!shippingOK) {
			Controller.dispatchEvent(new CustomEvent('addnotification', { detail: { type: 'error', content: createElement('span', { i18n: '#error_while_creating_order' }) } }));
			return;
		}

		this.#displayCheckout();
	}

	async #paymentComplete() {
		if (!this.#paymentCompleteDetails) {
			this.#navigateTo('/@checkout');
			return;
		}

		if (TESTING) {
			console.log(this.#paymentCompleteDetails);
		}

		this.#displayPaymentComplete();
		this.#broadcastChannel.postMessage({ action: BroadcastMessage.ReloadCart });
	}

	async #initProductPage(productId: string) {
		const product = await getShopProduct(productId);
		if (product) {
			this.#appContent.setProduct(product);
		} else {
			this.#navigateTo('/@products');
		}
	}

	async #initOrderPage(orderId: string) {
		const { requestId, response } = await fetchApi({
			action: 'get-order',
			version: 1,
			params: {
				order_id: orderId,
			},
		}) as { requestId: string, response: OrderResponse };
		if (response.success && response.result?.order) {
			const order = new Order();
			order.fromJSON(response.result.order);
			this.#appContent.setOrder(order);

		} else {
			Controller.dispatchEvent(new CustomEvent('addnotification', { detail: { type: 'error', content: createElement('span', { i18n: '#failed_to_get_order_details' }) } }));
		}
	}

	async #displayProducts() {
		const shopProducts = await this.#refreshProducts();

		if (shopProducts) {
			this.#appContent.setProducts(shopProducts);
		}
	}

	async #refreshProducts() {
		const { requestId, response } = await fetchApi({
			action: 'get-products',
			version: 1,
		}) as { requestId: string, response: GetProductsResponse };

		if (response?.success && response.result?.products) {
			console.log(response);
			const products: Array<Product> = [];
			for (const productJSON of response.result.products) {
				const product = new Product();
				product.fromJSON(productJSON);
				products.push(product);
			}
			return products;
		} else {
			//Controller.dispatchEvent(new CustomEvent('addnotification', {detail: {type: 'error', content: createElement('span', {i18n:'#error_while_sending_message'})}}));
		}
	}

	#initPage() {
		defineHarmonyCopy();
		defineHarmonySwitch();
		defineHarmonyPalette();
		defineHarmonySlideshow();
		createShadowRoot('div', {
			parent: document.body,
			adoptStyle: applicationCSS,
			childs: [
				this.#appToolbar.getHTML(),
				this.#appContent.getHTML(),
				this.#appFooter.getHTML(),
			]
		});
	}

	async #initSession() {
		//let response = await fetch('/getcurrency');
		//let json = await response.json();
		const result = await ServerAPI.getCurrency();
		//if (json && json.success) {
		this.#setCurrency(result);
		//}
	}

	#setCurrency(currency: string) {
		this.#appToolbar.setCurrency(/*currency*/);
	}

	#navigateTo(url: string, replaceSate = false) {
		history[replaceSate ? 'replaceState' : 'pushState']({}, '', url);
		this.#startup();
	}

	#scheduleRefreshProductPage() {
		clearTimeout(this.#refreshPageTimeout);
		this.#refreshPageTimeout = setTimeout(() => {
			if (this.#pageSubType == PAGE_SUBTYPE_SHOP_PRODUCT) {
				this.#navigateTo(document.location.pathname, true);
			}
		}, REFRESH_PRODUCT_PAGE_DELAY);
	}

	#pushState(url: string) {
		history.pushState({}, '', url);
	}

	#replaceState(url: string) {
		history.replaceState(this.#getHistoryState(), '', url);
	}

	#historyStateChanged() {
		history.replaceState(this.#getHistoryState(), '');
	}

	#getHistoryState() {
		return {
			//columnCartVisible: this.#htmlColumnCartVisible,
		};
	}

	#restoreHistoryState({ columnCartVisible = false } = {}) {
		//this.#htmlColumnCartVisible = columnCartVisible;
	}

	#onPaymentComplete(order: OrderJSON) {
		if (!this.#order) {
			this.#order = new Order();
		}
		this.#order.fromJSON(order);
		console.log(this.#order);
		this.#paymentCompleteDetails = { order: this.#order };
		this.#order = null;
		//this.#orderSummary.setOrder(null);
		this.#loadCart();

		this.#navigateTo(`/@order/${order.id}`);
	}
	/*
	set theme(theme) {
		document.documentElement.classList.remove('light');
		document.documentElement.classList.remove('dark');
		//document.documentElement.classList.add(theme);
	}
	*/

	async #processMessage(event: MessageEvent) {
		switch (event.data.action) {
			case BroadcastMessage.CartChanged:
				this.#cart.fromJSON(event.data.cart);
				const showColumnCart = this.#cart.totalQuantity > 0;
				//this.#htmlColumnCartVisible = showColumnCart;
				//this.#htmlColumnCart?.display(showColumnCart);
				this.#historyStateChanged();
				Controller.dispatchEvent(new CustomEvent(EVENT_REFRESH_CART, { detail: this.#cart }));
				Controller.dispatchEvent(new CustomEvent(EVENT_CART_COUNT, { detail: this.#cart.totalQuantity }));
				break;
			case BroadcastMessage.CartLoaded:
				Controller.dispatchEvent(new CustomEvent(EVENT_CART_COUNT, { detail: this.#cart.totalQuantity }));
				break;
			case BroadcastMessage.ReloadCart:
				this.#loadCart();
				break;
			case BroadcastMessage.FavoritesChanged:
				setFavorites(event.data.favorites);
				await this.#refreshFavorites();
				this.#countFavorites();
				break;
			case BroadcastMessage.FontSizeChanged:
				this.#setFontSize(event.data.fontSize);
				break;
		}
	}

	async #loadCart() {
		const { requestId, response } = await fetchApi({
			action: 'get-cart',
			version: 1,
		}) as { requestId: string, response: GetCartResponse };
		if (TESTING) {
			console.log(response);
		}

		if (response.success && response?.result?.cart) {
			this.#cart.fromJSON(response?.result?.cart);

			this.#refreshCart();

			this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartLoaded, cart: this.#cart.toJSON() });
		}
	}
}
new Application();
