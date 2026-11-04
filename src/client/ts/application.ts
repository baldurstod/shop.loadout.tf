import { addNotification, NotificationsPlacement, NotificationType, setNotificationsPlacement } from 'harmony-browser-utils';
import { themeCSS } from 'harmony-css';
import { createElement, createShadowRoot, defineHarmonyCopy, defineHarmonyPalette, defineHarmonySlideshow, defineHarmonySwitch, documentStyle, I18n } from 'harmony-ui';
import applicationCSS from '../css/application.css';
import htmlCSS from '../css/html.css';
import '../css/item.css';
import '../css/order.css';
import '../css/shop.css';
import '../css/vars.css';
import english from '../json/i18n/english.json';
import { setCurrency } from './appdatas';
import { BROADCAST_CHANNEL_NAME, PageSubType, PageType } from './constants';
import { AddToCartDetail, Controller, ControllerEvent, FavoriteDetail, LoginSuccessfulDetail, NavigateToDetail, SendContactDetail, SetQuantityDetail } from './controller';
import { RequestUserOrders, UserInfos } from './controllerevents';
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
import { GetOrdersResponse, InitCheckoutResponse, OrderJSON, OrderResponse, SetShippingAddressResponse, SetShippingMethodResponse } from './responses/order';
import { GetProductsResponse } from './responses/products';
import { addApiErrorNotification, addApiSuccessNotification } from './responses/response';
import { getShopProduct } from './shopproducts';
import { getUser, resetUser } from './user';
import { HTMLShopProductElement } from './view/components/shopproduct';
import { Footer } from './view/footer';
import { MainContent } from './view/maincontent';
import { Toolbar } from './view/toolbar';

const REFRESH_PRODUCT_PAGE_DELAY = 5000;

void documentStyle(htmlCSS);
void documentStyle(themeCSS);

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
	//#userResponse?: UserResponseResult;
	#redirect = '';

	constructor() {
		I18n.setOptions({ translations: [english] });
		I18n.start();
		setNotificationsPlacement(NotificationsPlacement.BottomRight);

		Controller.addEventListener(ControllerEvent.AddToCart, (event: Event) => { void this.#addToCart((event as CustomEvent<AddToCartDetail>).detail.productId, (event as CustomEvent<AddToCartDetail>).detail.quantity) });
		Controller.addEventListener(ControllerEvent.SetQuantity, (event: Event) => { void this.#setQuantity((event as CustomEvent<SetQuantityDetail>).detail.productId, (event as CustomEvent<SetQuantityDetail>).detail.quantity) });
		Controller.addEventListener(ControllerEvent.NavigateTo, (event: Event) => { void this.#navigateTo((event as CustomEvent<NavigateToDetail>).detail.url, (event as CustomEvent<NavigateToDetail>).detail.replaceSate) });
		//Controller.addEventListener('pushstate', (event: Event) => this.#pushState((event as CustomEvent).detail.url));
		//Controller.addEventListener('replacestate', (event: Event) => this.#replaceState((event as CustomEvent).detail.url));
		Controller.addEventListener(ControllerEvent.PaymentComplete, (event: Event) => { void this.#onPaymentComplete((event as CustomEvent<OrderJSON>).detail) });
		Controller.addEventListener(ControllerEvent.PaymentError, () => { void this.#onPaymentError() });
		Controller.addEventListener(ControllerEvent.Favorite, (event: Event) => { void this.#favorite((event as CustomEvent<FavoriteDetail>).detail.productId) });
		Controller.addEventListener(ControllerEvent.ScheduleRefreshProductPage, () => this.#scheduleRefreshProductPage());
		Controller.addEventListener(ControllerEvent.RefreshCart, () => { void this.#refreshCart() });
		Controller.addEventListener(ControllerEvent.UserInfoChanged, (event: Event) => this.#setUserInfos(event as CustomEvent<UserInfos>));
		Controller.addEventListener(ControllerEvent.RequestUserOrders, (event: Event) => { void this.#refreshUserOrders(event as CustomEvent<RequestUserOrders>) });
		Controller.addEventListener(ControllerEvent.PaymentCancelled, () => { void this.#paymentCancelled() });
		Controller.addEventListener(ControllerEvent.SetShippingMethod, () => { void this.#setShippingMethod() });

		Controller.addEventListener(ControllerEvent.LoginSuccessful, event => { void this.#onLoginSuccessful((event as CustomEvent<LoginSuccessfulDetail>).detail) });
		Controller.addEventListener(ControllerEvent.LogoutSuccessful, () => { void this.#onLogoutSuccessfull() });

		this.#initListeners();

		new BroadcastChannel(BROADCAST_CHANNEL_NAME).addEventListener('message', (event) => {
			void this.#processMessage(event);
		});

		//this.theme = 'light';
		void this.#init();
	}

	async #init(): Promise<void> {
		this.#initPage();
		await this.#initSession();
		await this.#loadCart();
		await this.#startup();
		await this.#initFavorites();
		await this.#initCountries();
		addEventListener('popstate', () => { void this.#startup(/*event.state ?? {}*/) });

	}

	#initListeners(): void {
		Controller.addEventListener(ControllerEvent.IncreaseFontSize, () => this.#changeFontSize(1));
		Controller.addEventListener(ControllerEvent.DecreaseFontSize, () => this.#changeFontSize(-1));
		Controller.addEventListener(ControllerEvent.SendContact, (event: Event) => { void this.#sendContact((event as CustomEvent<SendContactDetail>).detail) });
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
		const pathname = document.location.pathname;
		const user = await getUser();

		// If the user is not verified, disallow @checkout path and redirect to @verify
		if (user && !user.isEmailVerified() && pathname.includes('@checkout')) {
			await this.#navigateTo('/@verify');
			return;
		}

		this.#restoreHistoryState(/*historyState*/);
		this.#pageSubType = PageSubType.Unknown;
		switch (true) {
			case pathname.includes('@cart'):
				this.#pageType = PageType.Cart;
				break;
			case pathname.includes('@products'):
				this.#pageType = PageType.Products;
				await this.#displayProducts();
				break;
			case pathname.includes('@favorites'):
				this.#pageType = PageType.Favorites;
				break;
			case pathname.includes('@product'):
				this.#pageType = PageType.Product;
				await this.#initProductFromUrl();
				break;
			case pathname.includes('@checkout'):
				if (this.#cart.totalQuantity === 0) {
					await this.#navigateTo('/@products');
					break;
				}
				this.#pageType = PageType.Checkout;
				switch (document.location.hash) {
					case '':
						this.#pageSubType = PageSubType.CheckoutInit;
						await this.#initCheckout();
						break;
					case '#address':
						this.#pageSubType = PageSubType.CheckoutAddress;
						await this.#initAddress();
						break;
					case '#shipping':
						this.#pageSubType = PageSubType.CheckoutShipping;
						await this.#initShipping();
						break;
					case '#summary':
						this.#pageSubType = PageSubType.CheckoutSummary;
						await this.#initCheckoutSummary();
						break;
					case '#payment':
						this.#pageSubType = PageSubType.CheckoutPayment;
						await this.#initPayment();
						break;
					case '#complete':
						this.#pageSubType = PageSubType.CheckoutComplete;
						await this.#paymentComplete();
						break;
				}
				break;
			case pathname.includes('@login'):
				if (this.#authenticated) {
					await this.#navigateTo('/@user');
					return;
				}
				this.#pageType = PageType.Login;
				await this.#viewLoginPage();
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
				await this.#initOrderFromUrl();
				break;
			case pathname.includes('@user'):
				if (!this.#authenticated) {
					this.#redirect = '@user';
					await this.#navigateTo('/@login');
					return;
				}
				this.#pageType = PageType.User;
				break;
			case pathname.includes('@verify'):
				if (!this.#authenticated) {
					this.#redirect = '@user';
					await this.#navigateTo('/@login');
					return;
				}
				this.#pageType = PageType.Verify;
				break;
			default:
				await this.#navigateTo('/@products');
				break;
		}

		await this.#appContent.setActivePage(this.#pageType, this.#pageSubType);
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
		Controller.dispatchEvent<number>(ControllerEvent.FavoritesCount, { detail: favoritesCount() });
	}

	async #onLoginSuccessful(detail: LoginSuccessfulDetail): Promise<void> {
		addApiSuccessNotification('#login_successful');
		await this.#refreshUser();
		await this.#initFavorites();
		this.#broadcastChannel.postMessage({ action: BroadcastMessage.ReloadCart });

		const redirect = detail.redirect;
		if (redirect) {
			await this.#navigateTo(redirect);
		} else {
			if (this.#redirect == '') {
				await this.#navigateTo('/@products');
			} else {
				await this.#navigateTo(this.#redirect);
				this.#redirect = '';
			}
		}
	}

	async #onLogoutSuccessfull(): Promise<void> {
		addApiSuccessNotification('#logout_successful');
		this.#setAuthenticated(false);
		resetUser();
		await this.#navigateTo('/@products');
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

	async #refreshCart(): Promise<void> {
		await this.#appContent.setCart(this.#cart);
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

		await this.#refreshCart();
	}

	async #initOrderFromUrl(): Promise<void> {
		if (!this.#authenticated) {
			// Prevents from viewing order if not logged
			const notification = addNotification(createElement('span', {
				i18n: {
					innerText: '#you_must_login_to_use_this_functionality',
				},
			}), NotificationType.Info, 0);

			Controller.addEventListener(ControllerEvent.LoginSuccessful, function close(): void {
				notification.close();
				Controller.removeEventListener(ControllerEvent.LoginSuccessful, close, false);
			});

			this.#redirect = document.location.pathname;
			await this.#navigateTo('/@login');
			return;
		}

		const result = /@order\/([^\/]*)/i.exec(document.location.pathname);
		if (result) {
			await this.#loadCart();
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
			addApiSuccessNotification('#message_successfully_sent');
			//detail.callback(true);
		} else {
			addApiErrorNotification('#error_while_sending_message', requestId, response);
			Controller.dispatchEvent(ControllerEvent.SendContactError);
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

			Controller.addEventListener(ControllerEvent.LoginSuccessful, function close(): void {
				notification.close();
				Controller.removeEventListener(ControllerEvent.LoginSuccessful, close, false);
			});

			this.#redirect = '@checkout';
			await this.#navigateTo('/@login');
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

			await this.#navigateTo('/@checkout#address', true);

		} else {
			addApiErrorNotification('#failed_to_init_order', requestId, response);
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

	async #initAddress(): Promise<void> {
		if (!this.#order) {
			await this.#navigateTo('/@checkout');
			return;
		}

		this.#displayCheckout();
	}

	async #initShipping(): Promise<void> {
		if (!this.#order) {
			await this.#navigateTo('/@checkout');
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
			await this.#navigateTo('/@checkout');
			return false;
		}

		const { response, requestId } = await fetchApi('set-shipping-address', 1, {
			shipping_address: this.#order.shippingAddress,
			same_billing_address: this.#order.sameBillingAddress,
			...(!this.#order.sameBillingAddress && { billing_address: this.#order.billingAddress }),
		}) as { requestId: string, response: SetShippingAddressResponse };

		if (!response?.success) {
			addApiErrorNotification('#error_while_setting_shipping_address', requestId, response);
			Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@checkout#address' } });
			return false;
		}
		return true;
	}

	async #getShippingMethods(): Promise<boolean> {
		if (!this.#order) {
			await this.#navigateTo('/@checkout');
			return false;
		}

		const { requestId, response } = await fetchApi('get-shipping-methods', 1) as { requestId: string, response: SetShippingAddressResponse };

		if (response?.success && response.result?.order) {
			this.#order.fromJSON(response.result.order);
			this.#appContent.setCheckoutOrder(this.#order);
			return true;
		} else {
			addApiErrorNotification('#error_while_requesting_shipping_methods', requestId, response);

			Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@checkout#address' } });
			return false;
		}
	}

	async #sendShippingMethod(): Promise<{ requestId: string, shippingOK: boolean }> {
		if (!this.#order) {
			await this.#navigateTo('/@checkout');
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

		addApiErrorNotification('#error_while_creating_order', requestId, response);

		return { requestId: requestId, shippingOK: false };
	}

	async #initCheckoutSummary(): Promise<void> {
		if (!this.#order) {
			await this.#navigateTo('/@checkout');
			return;
		}
	}

	async #initPayment(): Promise<void> {
		if (!this.#order) {
			await this.#navigateTo('/@checkout');
			return;
		}

		this.#paymentCompleteDetails = null;

		/*
		const { requestId, shippingOK } = await this.#sendShippingMethod();
		if (!shippingOK) {
			return;
		}
		*/

		this.#displayCheckout();
	}

	async #setShippingMethod(): Promise<void> {
		if (!this.#order) {
			await this.#navigateTo('/@checkout');
			return;
		}

		const { shippingOK } = await this.#sendShippingMethod();
		if (shippingOK) {
			await this.#navigateTo('/@checkout#summary');
		}
	}

	async #paymentComplete(): Promise<void> {
		if (!this.#paymentCompleteDetails) {
			await this.#navigateTo('/@checkout');
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
			await this.#navigateTo('/@products');
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
			addApiErrorNotification('#failed_to_get_order_details', requestId, response);
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
		await this.#refreshUser();
	}

	async #refreshUser(): Promise<void> {
		const user = await getUser();
		if (user) {
			this.#setAuthenticated(true, user.getDisplayName());
			this.#appContent.refreshUserPage();
		}
	}

	async #navigateTo(url: string, replaceSate = false): Promise<void> {
		history[replaceSate ? 'replaceState' : 'pushState']({}, '', url);
		await this.#startup();
	}

	#scheduleRefreshProductPage(): void {
		const refresh = async (): Promise<void> => {
			if (this.#pageType == PageType.Product) {
				await this.#navigateTo(document.location.pathname, true);
			}
		}

		clearTimeout(this.#refreshPageTimeout);
		this.#refreshPageTimeout = setTimeout(() => { void refresh() }, REFRESH_PRODUCT_PAGE_DELAY);
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

	async #onPaymentComplete(order: OrderJSON): Promise<void> {
		if (!this.#order) {
			this.#order = new Order();
		}
		this.#order.fromJSON(order);
		console.log(this.#order);
		this.#paymentCompleteDetails = { order: this.#order };
		this.#order = null;
		//this.#orderSummary.setOrder(null);
		await this.#loadCart();
		this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartChanged, cart: this.#cart.toJSON() });

		await this.#navigateTo(`/@order/${order.id}`);
	}

	async #onPaymentError(): Promise<void> {
		await this.#navigateTo(`/@products`);
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
				Controller.dispatchEvent<Cart>(ControllerEvent.RefreshCart, { detail: this.#cart });
				Controller.dispatchEvent<number>(ControllerEvent.CartCount, { detail: this.#cart.totalQuantity });
				break;
			case BroadcastMessage.CartLoaded:
				Controller.dispatchEvent(ControllerEvent.CartCount, { detail: this.#cart.totalQuantity });
				break;
			case BroadcastMessage.ReloadCart:
				await this.#loadCart();
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

			await this.#refreshCart();

			this.#broadcastChannel.postMessage({ action: BroadcastMessage.CartLoaded, cart: this.#cart.toJSON() });
		}
	}

	#setAuthenticated(authenticated: boolean, displayName?: string): void {
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

	async #refreshUserOrders(event: CustomEvent<RequestUserOrders>): Promise<void> {
		const requestOrders = event.detail;

		const { response } = await fetchApi('get-orders', 1) as { requestId: string, response: GetOrdersResponse };
		if (response?.success) {
			const orders: Order[] = [];

			for (const orderJSON of response.result!.orders) {
				const o = new Order();
				o.fromJSON(orderJSON);
				orders.push(o);
			}

			requestOrders.callback(orders);
		} else {
			requestOrders.callback([]);
		}
	}

	async #paymentCancelled(/*event: CustomEvent<PaymentCancelled>*/): Promise<void> {
		//const paymentCancelled = event.detail;

		await this.#navigateTo('/@checkout#shipping');
	}
}
new Application();
