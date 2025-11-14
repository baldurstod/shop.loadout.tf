import { addNotification, NotificationsPlacement, NotificationType, setNotificationsPlacement } from 'harmony-browser-utils';
import { themeCSS } from 'harmony-css';
import { createElement, createShadowRoot, defineHarmonyCopy, defineHarmonyPalette, defineHarmonySlideshow, defineHarmonySwitch, documentStyle, I18n } from 'harmony-ui';
import { BROADCAST_CHANNEL_NAME, PageSubType, PageType } from './constants';
import { Controller } from './controller';
import { getShopProduct } from './shopproducts';
import { Footer } from './view/footer';
import { MainContent } from './view/maincontent';
import { Toolbar } from './view/toolbar';
//import { OrderSummary } from './view/ordersummary';
import applicationCSS from '../css/application.css';
import htmlCSS from '../css/html.css';
import '../css/item.css';
import '../css/order.css';
import '../css/shop.css';
import '../css/vars.css';
import english from '../json/i18n/english.json';
import { setCurrency } from './appdatas';
import { ControllerEvents, EVENT_CART_COUNT, EVENT_DECREASE_FONT_SIZE, EVENT_FAVORITES_COUNT, EVENT_INCREASE_FONT_SIZE, EVENT_NAVIGATE_TO, EVENT_REFRESH_CART, EVENT_SEND_CONTACT, EVENT_SEND_CONTACT_ERROR, RequestUserInfos, UserInfos } from './controllerevents';
import { BroadcastMessage, BroadcastMessageEvent, CartChangedEvent, FavoritesChangedEvent } from './enums';
import { favoritesCount, getFavorites, setFavorites, toggleFavorite } from './favorites';
import { fetchApi } from './fetchapi';
import { Cart } from './model/cart';
import { Countries } from './model/countries';
import { Order } from './model/order';
import { Product, setRetailPrice } from './model/product';
import { AddProductResponse, GetCartResponse } from './responses/cart';
import { CountriesResponse } from './responses/countries';
import { GetCurrencyResponse } from './responses/currency';
import { FavoritesResponse } from './responses/favorites';
import { InitCheckoutResponse, OrderJSON, OrderResponse, SetShippingAddressResponse, SetShippingMethodResponse } from './responses/order';
import { GetProductsResponse } from './responses/products';
import { GetUserResponse } from './responses/user';
import { HTMLShopProductElement } from './view/components/shopproduct';

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
	#authenticated = false;
	#displayName = '';
	#redirect = '';

	constructor() {
		I18n.setOptions({ translations: [english] });
		I18n.start();
		setNotificationsPlacement(NotificationsPlacement.BottomRight);

		Controller.addEventListener('addtocart', (event: Event) => { this.#addToCart((event as CustomEvent<{ product: string, quantity: number }>).detail.product, (event as CustomEvent<{ product: string, quantity: number }>).detail.quantity) });
		Controller.addEventListener('setquantity', (event: Event) => { this.#setQuantity((event as CustomEvent<{ id: string, quantity: number }>).detail.id, (event as CustomEvent<{ id: string, quantity: number }>).detail.quantity) });
		Controller.addEventListener(EVENT_NAVIGATE_TO, (event: Event) => this.#navigateTo((event as CustomEvent<{ url: string }>).detail.url, (event as CustomEvent<{ replaceSate: boolean }>).detail.replaceSate));
		//Controller.addEventListener('pushstate', (event: Event) => this.#pushState((event as CustomEvent).detail.url));
		//Controller.addEventListener('replacestate', (event: Event) => this.#replaceState((event as CustomEvent).detail.url));
		Controller.addEventListener('paymentcomplete', (event: Event) => this.#onPaymentComplete((event as CustomEvent<OrderJSON>).detail));
		Controller.addEventListener('favorite', (event: Event) => { this.#favorite((event as CustomEvent<{ productId: string }>).detail.productId) });
		Controller.addEventListener('schedulerefreshproductpage', () => this.#scheduleRefreshProductPage());
		Controller.addEventListener(EVENT_REFRESH_CART, () => this.#refreshCart());
		Controller.addEventListener(ControllerEvents.UserInfoChanged, (event: Event) => this.#setUserInfos(event as CustomEvent<UserInfos>));
		Controller.addEventListener(ControllerEvents.RequestUserInfos, (event: Event) => this.#requestUserInfos(event as CustomEvent<RequestUserInfos>));
		Controller.addEventListener(ControllerEvents.PaymentCancelled, () => this.#paymentCancelled(/*event as CustomEvent<PaymentCancelled>*/));

		Controller.addEventListener('loginsuccessful', (event: Event) => {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#login_successful',
				},
			}), NotificationType.Success, 4);
			this.setAuthenticated(true, (event as CustomEvent<{ displayName: string }>).detail.displayName);
			this.#initFavorites();
			this.#broadcastChannel.postMessage({ action: BroadcastMessage.ReloadCart });

			if (this.#redirect == '') {
				this.#navigateTo('/@products');
			} else {
				this.#navigateTo(this.#redirect);
				this.#redirect = '';
			}
		});
		Controller.addEventListener('logoutsuccessful', () => {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#logout_successful',
				},
			}), NotificationType.Success, 4);
			this.setAuthenticated(false);
			this.#navigateTo('/@products');
		});

		this.#initListeners();

		new BroadcastChannel(BROADCAST_CHANNEL_NAME).addEventListener('message', (event) => {
			this.#processMessage(event);
		});

		//this.theme = 'light';
		this.#init();
	}

	async #init(): Promise<void> {
		this.#initPage();
		await this.#initSession();
		await this.#startup();
		await this.#initFavorites();
		await this.#initCountries();
		addEventListener('popstate', () => { this.#startup(/*event.state ?? {}*/) });
		await this.#loadCart();

	}

	#initListeners(): void {
		Controller.addEventListener(EVENT_INCREASE_FONT_SIZE, () => this.#changeFontSize(1));
		Controller.addEventListener(EVENT_DECREASE_FONT_SIZE, () => this.#changeFontSize(-1));
		Controller.addEventListener(EVENT_SEND_CONTACT, (event: Event) => { this.#sendContact((event as CustomEvent<{ subject: string, email: string, content: string, }>).detail) });
	}

	#changeFontSize(change: number): void {
		let size = this.#fontSize;
		if (change > 0) {
			size *= 1.1;
		} else {
			size /= 1.1;
		}
		this.#broadcastChannel.postMessage({ action: BroadcastMessage.FontSizeChanged, fontSize: size });
	}

	#setFontSize(size: number): void {
		this.#fontSize = size;
		document.documentElement.style.setProperty('--font-size', `${this.#fontSize}px`);
	}

	async #startup(/*historyState = {}*/): Promise<void> {
		this.#restoreHistoryState(/*historyState*/);
		const pathname = document.location.pathname;
		this.#pageSubType = PageSubType.Unknown;
		switch (true) {
			case pathname.includes('@cart'):
				this.#pageType = PageType.Cart;
				break;
			case pathname.includes('@products'):
				this.#pageType = PageType.Products;
				this.#displayProducts();
				break;
			case pathname.includes('@favorites'):
				this.#pageType = PageType.Favorites;
				break;
			case pathname.includes('@product'):
				this.#pageType = PageType.Product;
				await this.#initProductFromUrl();
				break;
			case pathname.includes('@checkout'):
				this.#pageType = PageType.Checkout;
				switch (document.location.hash) {
					case '':
						this.#pageSubType = PageSubType.CheckoutInit;
						this.#initCheckout();
						break;
					case '#address':
						this.#pageSubType = PageSubType.CheckoutAddress;
						this.#initAddress();
						break;
					case '#shipping':
						this.#pageSubType = PageSubType.CheckoutShipping;
						this.#initShipping();
						break;
					case '#payment':
						this.#pageSubType = PageSubType.CheckoutPayment;
						await this.#initPayment();
						break;
					case '#complete':
						this.#pageSubType = PageSubType.CheckoutComplete;
						this.#paymentComplete();
						break;
				}
				break;
			case pathname.includes('@login'):
				if (this.#authenticated) {
					this.#navigateTo('/@user');
					return;
				}
				this.#pageType = PageType.Login;
				this.#viewLoginPage();
				break;
			case pathname.includes('@cookies'):
				this.#pageType = PageType.Cookies;
				break;
			case pathname.includes('@privacy'):
				this.#pageType = PageType.Privacy;
				break;
			case pathname.includes('@contact'):
				this.#pageType = PageType.Contact;
				break;
			case pathname.includes('@order'):
				this.#pageType = PageType.Order;
				this.#initOrderFromUrl();
				break;
			case pathname.includes('@user'):
				if (!this.#authenticated) {
					this.#navigateTo('/@login');
					return;
				}
				this.#pageType = PageType.User;
				break;
			default:
				this.#navigateTo('/@products');
				break;
		}

		this.#appContent.setActivePage(this.#pageType, this.#pageSubType);
	}

	async #initFavorites(): Promise<void> {
		const { response } = await fetchApi('get-favorites', 1) as { requestId: string, response: FavoritesResponse };
		if (response?.success) {
			setFavorites(response.result?.favorites);

			this.#countFavorites();
			this.#broadcastChannel.postMessage({ action: BroadcastMessage.FavoritesChanged, favorites: getFavorites() });
		}
	}

	async #initCountries(): Promise<void> {
		const { response } = await fetchApi('get-countries', 1) as { requestId: string, response: CountriesResponse };
		if (response?.success && response.result?.countries) {
			this.#countries.fromJSON(response.result.countries);
			this.#appContent.setCountries(this.#countries);
		}
	}

	async #favorite(productId: string): Promise<void> {
		await fetchApi('set-favorite', 1, {
			product_id: productId,
			is_favorite: toggleFavorite(productId),
		});

		this.#broadcastChannel.postMessage({ action: BroadcastMessage.FavoritesChanged, favorites: getFavorites() });
		this.#countFavorites();
	}

	#countFavorites(): void {
		Controller.dispatchEvent(new CustomEvent(EVENT_FAVORITES_COUNT, { detail: favoritesCount() }));
	}

	async #addToCart(productId: string, quantity = 1): Promise<void> {
		if (TESTING) {
			console.log(productId, quantity);
		}

		const { response } = await fetchApi('add-product', 1, {
			product_id: productId,
			quantity: quantity,
		}) as { requestId: string, response: AddProductResponse };

		if (response.success && response.result?.cart) {
			this.#cart.fromJSON(response.result.cart);
			this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartChanged, cart: this.#cart.toJSON() });
		}

	}

	async #setQuantity(productId: string, quantity = 1): Promise<void> {
		if (TESTING) {
			console.log(productId, quantity);
		}

		const { response } = await fetchApi('set-product-quantity', 1, {
			product_id: productId,
			quantity: quantity,
		}) as { requestId: string, response: AddProductResponse };

		if (response.success && response.result?.cart) {
			this.#cart.fromJSON(response.result.cart);
			this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartChanged, cart: this.#cart.toJSON() });
		}
	}

	#refreshCart(): void {
		this.#appContent.setCart(this.#cart);
	}

	async #refreshFavorites(): Promise<void> {
		const favorites: Product[] = [];

		for (const productID of getFavorites()) {
			const product = await getShopProduct(productID);
			if (product) {
				favorites.push(product);
			}
		}

		this.#appContent.setFavorites(favorites);
	}

	async #initProductFromUrl(): Promise<void> {
		const result = /@product\/([^\/]*)/i.exec(document.location.pathname);
		if (result) {
			await this.#initProductPage(result[1]!);
		}

		this.#refreshCart();
	}

	async #initOrderFromUrl(): Promise<void> {
		const result = /@order\/([^\/]*)/i.exec(document.location.pathname);
		if (result) {
			this.#loadCart();
			await this.#initOrderPage(result[1]!);
		}
	}

	async #viewLoginPage(): Promise<void> {
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

	async #sendContact(detail: { subject: string, email: string, content: string, }): Promise<void> {
		const { requestId, response } = await fetchApi('send-message', 1, {
			subject: detail.subject,
			email: detail.email,
			content: detail.content,
		});


		if (response?.success) {
			addNotification(createElement('span', { i18n: '#message_successfully_sent' }), NotificationType.Success, 4);
			//detail.callback(true);
		} else {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_while_sending_message',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
			Controller.dispatchEvent(new CustomEvent(EVENT_SEND_CONTACT_ERROR));
		}
	}

	#displayCheckout(): void {
		return;
		/*
		let htmlCheckoutPage = createElement('div', { class: 'shop-checkout-page' });

		htmlCheckoutPage.append(this.#orderSummary.html);
		this.htmlContent.append(htmlCheckoutPage);
		this.#refreshOrder();
		*/
	}

	#displayPaymentComplete(): void {
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
							innerText: order.id,
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

	async #initCheckout(): Promise<void> {
		if (!this.#authenticated) {
			const notification = addNotification(createElement('span', {
				i18n: {
					innerText: '#you_must_login_to_use_this_functionality',
				},
			}), NotificationType.Info, 0);

			Controller.addEventListener('loginsuccessful', function close(): void {
				notification.close();
				Controller.removeEventListener('loginsuccessful', close, false);
			});

			this.#redirect = '@checkout';
			this.#navigateTo('/@login');
			return;
		}

		const { requestId, response } = await fetchApi('init-checkout', 1) as { requestId: string, response: InitCheckoutResponse };
		if (response?.success && response.result?.order) {
			const order = new Order();
			order.fromJSON(response.result.order);
			this.#appContent.setCheckoutOrder(order);

			this.#order = order;

			//this.#orderSummary.setOrder(order);
			this.#orderId = order.id;

			this.#navigateTo('/@checkout#address', true);

		} else {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#failed_to_init_order',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
			return;
		}
	}

	/*
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

		const json = await response.json();
		if (json?.success) {
			const order = new Order();
			order.fromJSON(json.result);
			//this.#orderSummary.setOrder(order);
		}
	}
	*/

	#initAddress(): void {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}

		this.#displayCheckout();
	}

	async #initShipping(): Promise<void> {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}

		let shippingOk = await this.#sendShippingAddress();
		shippingOk = shippingOk && await this.#getShippingMethods();
		if (shippingOk) {
			this.#displayCheckout();
		}
	}

	async #sendShippingAddress(): Promise<boolean> {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return false;
		}

		const { response } = await fetchApi('set-shipping-address', 1, {
			shipping_address: this.#order.shippingAddress,
			same_billing_address: this.#order.sameBillingAddress,
			...(!this.#order.sameBillingAddress && { billing_address: this.#order.billingAddress }),
		}) as { requestId: string, response: SetShippingAddressResponse };

		if (!response?.success) {
			addNotification(createElement('span', { innerText: response.error }), NotificationType.Error, 0);
			Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout#address' } }));
			return false;
		}
		return true;
	}

	async #getShippingMethods(): Promise<boolean> {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return false;
		}

		const { requestId, response } = await fetchApi('get-shipping-methods', 1) as { requestId: string, response: SetShippingAddressResponse };

		if (response?.success && response.result?.order) {
			this.#order.fromJSON(response.result.order);
			this.#appContent.setCheckoutOrder(this.#order);
			return true;
		} else {
			//addNotification(createElement('span', { innerText: response.error }), NotificationType.Error, 0);
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_request_id',
					values: {
						error: response.error,
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);

			Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout#address' } }));
			return false;
		}
	}

	async #sendShippingMethod(): Promise<{ requestId: string, shippingOK: boolean }> {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return { requestId: '0', shippingOK: false };
		}

		const { requestId, response } = await fetchApi('set-shipping-method', 1, {
			method: this.#order.shippingMethod,
		}) as { requestId: string, response: SetShippingMethodResponse };

		if (response?.success && response.result?.order) {
			this.#order.fromJSON(response.result.order);
			this.#appContent.setCheckoutOrder(this.#order);
			return { requestId: requestId, shippingOK: true };
		}
		return { requestId: requestId, shippingOK: false };
	}

	async #initPayment(): Promise<void> {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}

		this.#paymentCompleteDetails = null;

		const { requestId, shippingOK } = await this.#sendShippingMethod();
		if (!shippingOK) {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#error_while_creating_order',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
			return;
		}

		this.#displayCheckout();
	}

	#paymentComplete(): void {
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

	async #initProductPage(productId: string): Promise<void> {
		const product = await getShopProduct(productId);
		if (product) {
			this.#appContent.setProduct(product);
		} else {
			this.#navigateTo('/@products');
		}
	}

	async #initOrderPage(orderId: string): Promise<void> {
		const { requestId, response } = await fetchApi('get-order', 1, {
			order_id: orderId,
		}) as { requestId: string, response: OrderResponse };
		if (response.success && response.result?.order) {
			const order = new Order();
			order.fromJSON(response.result.order);
			this.#appContent.setOrder(order);

		} else {
			addNotification(createElement('span', {
				i18n: {
					innerText: '#failed_to_get_order_details',
					values: {
						requestId: requestId,
					},
				},
			}), NotificationType.Error, 0);
		}
	}

	async #displayProducts(): Promise<void> {
		const shopProducts = await this.#refreshProducts();

		if (shopProducts) {
			this.#appContent.setProducts(shopProducts);
		}
	}

	async #refreshProducts(): Promise<Product[] | null> {
		const { response } = await fetchApi('get-products', 1) as { requestId: string, response: GetProductsResponse };

		if (response?.success && response.result?.products) {
			console.log(response);
			const products: Product[] = [];
			for (const productJSON of response.result.products) {
				const product = new Product();
				product.fromJSON(productJSON);
				products.push(product);
			}

			const prices = response.result?.prices
			if (prices) {
				const currency = prices.currency;
				for (const productID in prices.prices) {
					setRetailPrice(currency, productID, prices.prices[productID]!);
				}
			}

			return products;
		} else {
			//Controller.dispatchEvent(new CustomEvent('addnotification', {detail: {type: 'error', content: createElement('span', {i18n:'#error_while_sending_message'})}}));
		}
		return null;
	}

	#initPage(): void {
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

	async #initSession(): Promise<void> {
		const { response } = await fetchApi('get-currency', 1) as { requestId: string, response: GetCurrencyResponse };
		if (response.success) {
			setCurrency(response.result!.currency);
		}

		const { response: orderResponse } = await fetchApi('get-active-order', 1) as { requestId: string, response: OrderResponse };
		if (orderResponse.success && orderResponse.result!.order) {
			this.#order = new Order();
			this.#order.fromJSON(orderResponse.result!.order);
			this.#appContent.setCheckoutOrder(this.#order);
		}

		const { response: userResponse } = await fetchApi('get-user', 1) as { requestId: string, response: GetUserResponse };
		if (response.success) {
			setCurrency(response.result!.currency);
			this.setAuthenticated(userResponse.result!.authenticated, userResponse.result!.display_name);
		}
	}

	#navigateTo(url: string, replaceSate = false): void {
		history[replaceSate ? 'replaceState' : 'pushState']({}, '', url);
		this.#startup();
	}

	#scheduleRefreshProductPage(): void {
		clearTimeout(this.#refreshPageTimeout);
		this.#refreshPageTimeout = setTimeout(() => {
			if (this.#pageSubType == PageSubType.ShopProduct) {
				this.#navigateTo(document.location.pathname, true);
			}
		}, REFRESH_PRODUCT_PAGE_DELAY);
	}

	#pushState(url: string): void {
		history.pushState({}, '', url);
	}

	#replaceState(url: string): void {
		history.replaceState(this.#getHistoryState(), '', url);
	}

	#historyStateChanged(): void {
		history.replaceState(this.#getHistoryState(), '');
	}

	#getHistoryState(): object {
		return {
			//columnCartVisible: this.#htmlColumnCartVisible,
		};
	}

	#restoreHistoryState(/*{ columnCartVisible = false } = {}*/): void {
		//this.#htmlColumnCartVisible = columnCartVisible;
	}

	#onPaymentComplete(order: OrderJSON): void {
		if (!this.#order) {
			this.#order = new Order();
		}
		this.#order.fromJSON(order);
		console.log(this.#order);
		this.#paymentCompleteDetails = { order: this.#order };
		this.#order = null;
		//this.#orderSummary.setOrder(null);
		this.#loadCart();
		this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartChanged, cart: this.#cart.toJSON() });

		this.#navigateTo(`/@order/${order.id}`);
	}
	/*
	set theme(theme) {
		document.documentElement.classList.remove('light');
		document.documentElement.classList.remove('dark');
		//document.documentElement.classList.add(theme);
	}
	*/

	async #processMessage(event: MessageEvent): Promise<void> {
		switch ((event as MessageEvent<BroadcastMessageEvent>).data.action) {
			case BroadcastMessage.CartChanged:
				this.#cart.fromJSON((event as MessageEvent<CartChangedEvent>).data.cart);
				//const showColumnCart = this.#cart.totalQuantity > 0;
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
				setFavorites((event as MessageEvent<FavoritesChangedEvent>).data.favorites);
				await this.#refreshFavorites();
				this.#countFavorites();
				break;
			case BroadcastMessage.FontSizeChanged:
				this.#setFontSize((event as MessageEvent<{ fontSize: number }>).data.fontSize);
				break;
		}
	}

	async #loadCart(): Promise<void> {
		const { response } = await fetchApi('get-cart', 1) as { requestId: string, response: GetCartResponse };
		if (TESTING) {
			console.log(response);
		}

		if (response.success && response?.result?.cart) {
			this.#cart.fromJSON(response?.result?.cart);

			this.#refreshCart();

			this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartLoaded, cart: this.#cart.toJSON() });
		}
	}

	setAuthenticated(authenticated: boolean, displayName?: string): void {
		this.#authenticated = authenticated;
		this.#appToolbar.setAuthenticated(authenticated);
		this.#displayName = displayName ?? '';
		this.#appToolbar.setDisplayName(this.#displayName);
	}

	#setUserInfos(event: CustomEvent<UserInfos>): void {
		const userInfos = event.detail;

		if (userInfos.displayName !== undefined) {
			this.#appToolbar.setDisplayName(userInfos.displayName);
		}
	}

	#requestUserInfos(event: CustomEvent<RequestUserInfos>): void {
		const requestUserInfos = event.detail;

		requestUserInfos.callback({
			authenticated: this.#authenticated,
			displayName: this.#displayName,
		})
	}

	#paymentCancelled(/*event: CustomEvent<PaymentCancelled>*/): void {
		//const paymentCancelled = event.detail;

		this.#navigateTo('/@checkout#shipping');
	}
}
new Application();
