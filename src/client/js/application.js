import {NotificationManager} from 'harmony-browser-utils/src/NotificationManager.js';
import { themeCSS } from 'harmony-css';
import { bookmarksPlainSVG, shoppingCartSVG, textIncreaseSVG, textDecreaseSVG } from 'harmony-svg';
import {createElement, hide, show, display, I18n, documentStyle} from 'harmony-ui';

import { getCartTotalPriceFormatted } from './carttotalprice.js';
import { getShopProduct } from './shopproducts.js';
import {PAYPAL_APP_CLIENT_ID, ALLOWED_CURRENCIES, BROADCAST_CHANNEL_NAME, PAGE_TYPE_CART, PAGE_TYPE_CHECKOUT, PAGE_TYPE_PRODUCTS, PAGE_TYPE_COOKIES, PAGE_TYPE_PRIVACY, PAGE_TYPE_CONTACT, PAGE_TYPE_LOGIN, PAGE_TYPE_ORDER, PAGE_TYPE_PRODUCT, PAGE_TYPE_FAVORITES, PAGE_SUBTYPE_CHECKOUT_INIT, PAGE_SUBTYPE_CHECKOUT_ADDRESS, PAGE_SUBTYPE_CHECKOUT_SHIPPING, PAGE_SUBTYPE_CHECKOUT_PAYMENT, PAGE_SUBTYPE_CHECKOUT_COMPLETE } from './constants.js';
import {Controller} from './controller.js';
import {formatPrice, loadScript, formatPriceRange} from './utils.js';
import { Footer } from './view/footer.js';
import { MainContent } from './view/maincontent.js';
import { Toolbar } from './view/toolbar.js';
import {OrderView} from './view/orderview.js';
import {OrderSummary} from './view/ordersummary.js';
import { orderSummary } from './view/ordersummary2.js';
import {transactionSummary} from './view/transactionsummary.js';
import {CartProducts} from './view/components/cartproducts.js';
import {ColumnCart} from './view/components/columncart.js';
//import { ShopProductElement } from './view/components/shopproductelement.js';
//import {LabelProperty} from './view/components/labelproperty.js';
import {Cart} from './model/cart.js';
import {File} from './model/file.js';
import {Order} from './model/order.js';
import {OrderItem} from './model/orderitem.js';
import { ShopProduct } from './model/shopproduct.js';
import {SyncProduct} from './model/syncproduct.js';
import {SyncVariant} from './model/syncvariant.js';

import 'harmony-ui/dist/define/harmony-label-property.js';
import 'harmony-ui/dist/define/harmony-copy.js';

import '../css/header.css';
import '../css/item.css';
import '../css/order.css';
import '../css/shop.css';
import '../css/vars.css';

import applicationCSS from '../css/application.css';
import htmlCSS from '../css/html.css';

import english from '../json/i18n/english.json';

import { fetchApi } from './fetchapi.js';
import { ServerAPI } from './serverapi.js';
import { EVENT_CART_COUNT, EVENT_DECREASE_FONT_SIZE, EVENT_FAVORITES_COUNT, EVENT_INCREASE_FONT_SIZE, EVENT_NAVIGATE_TO, EVENT_REFRESH_CART, EVENT_SEND_CONTACT, EVENT_SEND_CONTACT_ERROR } from './controllerevents.js';
import { Countries } from './model/countries.js';

const REFRESH_PRODUCT_PAGE_DELAY = 20000;

documentStyle(htmlCSS);
documentStyle(themeCSS);


class Application {
	#htmlElement;
	#appToolbar = new Toolbar();
	#appContent = new MainContent();
	#appFooter = new Footer();
	#pageType;
	#pageSubType;
	#htmlFooter;
	#htmlContactPage;
	#htmlPrivacyPage;
	#htmlCookiesPage;
	#htmlFavoriteList;
	#htmlContactSubject;
	#htmlContactEmail;
	#htmlContactContent;
	#htmlContactButton;
	#order;
	//#orderView;
	#orderSummary;
	#currency;
	#htmlCurrency;
	#printfulOrder;
	#favorites = [];
	#broadcastChannel = new BroadcastChannel(BROADCAST_CHANNEL_NAME);
	#paymentCompleteDetails;
	#htmlColumnCart;
	#htmlColumnCartVisible = false;
	#htmlProductsPage;
	#htmlShopProduct;
	#htmlCartList;
	#htmlCheckoutButton;
	#htmlCart;
	#htmlHeaderFavorites;
	#syncProducts = new Map();
	#syncVariants = new Map();
	#orderId;
	#fontSize = 16;
	#refreshPageTimeout;
	#cart = new Cart();
	#countries = new Countries();
	constructor() {
		this.page;
//			this.#order = new Order();
		//this.#orderView = new OrderView();
		this.#orderSummary = new OrderSummary();
		//this.#order.cart = this.#cart;
		I18n.setOptions({translations:[english]});
		I18n.start();

		Controller.addEventListener('addtocart', (event) => this.#addToCart(event.detail.product, event.detail.quantity));
		Controller.addEventListener('setquantity', (event) => this.#setQuantity(event.detail.id, event.detail.quantity));
		Controller.addEventListener(EVENT_NAVIGATE_TO, (event) => this.#navigateTo(event.detail.url, event.detail.replaceSate));
		Controller.addEventListener('pushstate', (event) => this.#pushState(event.detail.url));
		Controller.addEventListener('replacestate', (event) => this.#replaceState(event.detail.url));
		Controller.addEventListener('addnotification', (event) => NotificationManager.addNotification(event.detail.content, event.detail.type));
		Controller.addEventListener('paymentcomplete', (event) => this.#onPaymentComplete(event.detail));
		Controller.addEventListener('favorite', (event) => this.#favorite(event.detail.productId));
		Controller.addEventListener('schedulerefreshproductpage', (event) => this.#scheduleRefreshProductPage());
		Controller.addEventListener(EVENT_REFRESH_CART, () => this.#refreshCart());
		this.#initListeners();

		new BroadcastChannel(BROADCAST_CHANNEL_NAME).addEventListener('message', (event) => {
			this.#processMessage(event);
		});

		//this.#cart.addEventListener('changed', (event) => this.#refreshCart());
		this.theme = 'light';

		//this.#cart.loadFromLocalStorage();
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
		Controller.addEventListener(EVENT_SEND_CONTACT, event => this.#sendContact(event.detail));
	}

	#changeFontSize(change) {
		if (change > 0) {
			this.#fontSize *= 1.1;
		} else {
			this.#fontSize /= 1.1;

		}
		document.documentElement.style.setProperty('--font-size', `${this.#fontSize}px`);

	}


	async #startup(historyState = {}) {
		this.#restoreHistoryState(historyState);
		let pathname = document.location.pathname;
		this.#pageSubType = null;
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
				//this.#viewCookiesPage();
				break;
			case pathname.includes('@privacy'):
				this.#pageType = PAGE_TYPE_PRIVACY;
				//this.#viewPrivacyPage();
				break;
			case pathname.includes('@contact'):
				this.#pageType = PAGE_TYPE_CONTACT;
				//this.#viewContactPage();
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
		});
		if (response?.success) {
			this.#favorites = response.result ?? [];
			this.#countFavorites();
			this.#broadcastChannel.postMessage({action: 'favoriteschanged', favorites: this.#favorites});
		}
	}

	async #initCountries() {
		const { requestId, response } = await fetchApi({
			action: 'get-countries',
			version: 1,
		});
		if (response?.success) {
			this.#countries.fromJSON(response.result.countries);
			this.#appContent.setCountries(this.#countries);
		}
	}

	async #favorite(productId) {
		const index = this.#favorites.indexOf(productId);
		let favorite = index > -1;//this.#favorites[productId];
		if (favorite) {
			//delete this.#favorites[productId];
			this.#favorites.splice(index, 1);
		} else {
			this.#favorites.push(productId);
			//this.#favorites[productId] = 1;
		}

		/*await fetch('/api', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({
				action: 'set-favorite',
				version: 1,
				params: {
					productId: productId,
					isFavorite: !favorite,
				}
			}),
		});*/

		await fetchApi({
			action: 'set-favorite',
			version: 1,
			params: {
				product_id: productId,
				is_favorite: !favorite,
			},
		});


		this.#broadcastChannel.postMessage({ action: 'favoriteschanged', favorites: this.#favorites });
		this.#countFavorites();
	}

	#countFavorites() {
		let count = 0;
		let favorites = this.#favorites;
		for (let externalProductId in favorites) {
			++count;
		}
		//this.#htmlHeaderFavorites.innerHTML = count;
		Controller.dispatchEvent(new CustomEvent(EVENT_FAVORITES_COUNT, { detail: count }));
	}

	async #addToCart(productId, quantity = 1) {
		if (TESTING) {
			console.log(productId, quantity);
		}

		//this.#cart.addProduct(productId, quantity);

		const { requestId, response } = await fetchApi({
			action: 'add-product',
			version: 1,
			params: {
				product_id: productId,
				quantity: quantity,
				//cart: this.#cart.toJSON(),
			},
		});

		//this.#saveCart();

		if (response.success) {
			this.#cart.fromJSON(response.result.cart);
			this.#broadcastChannel.postMessage({ action: 'cartchanged', cart: this.#cart.toJSON() });
		}

	}

	async #setQuantity(productId, quantity = 1) {
		if (TESTING) {
			console.log(productId, quantity);
		}

		const { requestId, response } = await fetchApi({
			action: 'set-product-quantity',
			version: 1,
			params: {
				product_id: productId,
				quantity: quantity,
				//cart: this.#cart.toJSON(),
			},
		});

		//this.#cart.setQuantity(productId, quantity);

		//this.#saveCart();
		//this.#broadcastChannel.postMessage({ action: 'cartchanged', cart: this.#cart.toJSON() });

		if (response.success) {
			this.#cart.fromJSON(response.result.cart);
			this.#broadcastChannel.postMessage({ action: 'cartchanged', cart: this.#cart.toJSON() });
		}
	}

	async #saveCart() {
		/*const cartJSON = [];
		for (let [productId, product] of this.#cart.items) {
			cartJSON.push({ productId: productId, quantity: product.quantity });
		}*/

		/*let response = await fetch('/setcart/', {
			method: 'POST',
			headers: {
				'Content-Type': 'application/json',
			},
			body: JSON.stringify({ cart: this.#cart.toJSON() }),
		});*/
		const { requestId, response } = await fetchApi({
			action: 'set-cart',
			version: 1,
			params: {
				cart: this.#cart.toJSON(),
			},
		});

		console.info(response);

		if (response?.success) {
			this.#broadcastChannel.postMessage({action: 'reloadcart'});
		}
	}

	async #refreshCart() {
		this.#appContent.setCart(this.#cart);
	}

	async #refreshFavorites() {
		const favorites = [];

		for (const productID of this.#favorites) {
			const product = await getShopProduct(productID);
			favorites.push(product)
		}


		this.#appContent.setFavorites(favorites);
		return;

		if (this.#htmlFavoriteList) {
			this.#htmlFavoriteList.innerText = '';
			let favorites = this.#favorites;
			let count = 0;

			for (const shopProduct of favorites) {
				//const isFavorited = favorites[productId];
				//if (isFavorited) {
				++count;
				await this.#createHtmlFavorite(shopProduct, this.#htmlFavoriteList);
				//}
				/*for (let externalVariantId in product) {
					let variant = product[externalVariantId];
					//console.log(externalProductId);
					++count;
					await this.#createHtmlFavorite(externalProductId, externalVariantId, this.#htmlFavoriteList);
				}*/
			}

			if (!count) {
				createElement('div', {
					parent: this.#htmlFavoriteList,
					i18n: '#empty_favorites_list',
				});

			}
		}
	}

	async #createHtmlFavorite(productId, htmlFavoriteList) {
		const shopProduct = await getShopProduct(productId);
		if (!shopProduct) {
			return;
		}

		const link = `/@product/${productId}`;

		createElement('div', {
			class: 'shop-product',
			parent: htmlFavoriteList,
			events: {
				click: () => this.#navigateTo(link),
				mouseup: (event) => {
					if (event.button == 1) {
						open(link, '_blank');
					}
				},
			},
			childs: [
				createElement('img', {
					class: 'shop-product-thumb',
					src: shopProduct.thumbnailUrl,//getThumbnailUrl('preview') ?? shopProduct.getThumbnailUrl('default'),
				}),
				createElement('div', {
					class: 'shop-product-description',
					childs: [
						createElement('div', {
							class: 'shop-product-name',
							innerHTML: shopProduct.name
						}),
						createElement('div', {
							class: 'shop-product-price',
							innerHTML: formatPrice(shopProduct.retailPrice)
						}),
					]
				}),
			]
		});
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
		let htmlLoginPage = createElement('div', {class:'shop-login-page'});
		htmlLoginPage.append(createElement('div', {id:'paypal-login-container'}));
		await loadScript('https://www.paypalobjects.com/js/external/api.js');

		paypal.use( ['login'], function (login) {
			login.render ({
				"appid":PAYPAL_APP_CLIENT_ID,
				"authend":"sandbox",
				"scopes":"profile email address",
				"containerid":"paypal-login-container",
				"responseType":"code",
				"locale":"en-us",
				"buttonType":"LWP",
				"buttonShape":"pill",
				"buttonSize":"lg",
				"fullPage":"true",
				"returnurl":document.location.origin + '/paypalredirect',
				"nonce":"12345678"
			});
		});

		this.htmlContent.append(htmlLoginPage);
	}

	/*
	async #viewCookiesPage() {
		this.#htmlCookiesPage = createElement('div', {
			class:'shop-cookies-page',
			parent: this.htmlContent,
			i18n: '#cookies_policy_content',
		});
	}
	*/

	/*async #viewPrivacyPage() {
		this.#htmlPrivacyPage = createElement('div', {
			class:'shop-privacy-page',
			parent: this.htmlContent,
			i18n: '#privacy_policy_content',
		});
	}*/

	async #viewContactPage() {
		this.#htmlContactPage = createElement('div', {
			class: 'shop-contact-page',
			parent: this.htmlContent,
			childs: [
				createElement('h1', {
					i18n: '#contact',
				}),
				createElement('div', {
					class: 'shop-contact-page-content',
					childs: [
						createElement('div', {i18n: '#subject'}),
						this.#htmlContactSubject = createElement('input'),
						createElement('div', {i18n: '#email'}),
						this.#htmlContactEmail = createElement('input'),
						createElement('div', {
							i18n: '#content',
						}),
						this.#htmlContactContent = createElement('textarea', {
							rows:10,
							cols:80,
						}),
						createElement('div', {
							childs: [
								this.#htmlContactButton = createElement('button', {
									i18n: '#send',
									events: {
										click: () => this.#sendContact(),
									},
								}),
							]
						}),
					]
				}),
			]
		});
	}

	async #sendContact(detail) {
		//this.#htmlContactButton.disabled = true;

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
			Controller.dispatchEvent(new CustomEvent('addnotification', {detail: {type: 'success', content: createElement('span', {i18n:'#message_successfully_sent'})}}));
			//detail.callback(true);
		} else {
			Controller.dispatchEvent(new CustomEvent('addnotification', {detail: {type: 'error', content: createElement('span', {i18n:'#error_while_sending_message'})}}));
			//detail.callback(false);
			//this.#htmlContactButton.disabled = false;
			Controller.dispatchEvent(new CustomEvent(EVENT_SEND_CONTACT_ERROR));
		}
	}

	#displayCheckout() {
		return;
		let htmlCheckoutPage = createElement('div', {class:'shop-checkout-page'});

		//htmlCheckoutPage.append(this.#orderView.html);
		htmlCheckoutPage.append(this.#orderSummary.html);
		this.htmlContent.append(htmlCheckoutPage);
		this.#refreshOrder();
	}

	#displayPaymentComplete() {
		let cartItems = [];

		let order = this.#paymentCompleteDetails.order;
		let paymentDetail = this.#paymentCompleteDetails.paymentDetail;
		if (TESTING) {
			console.log(order, paymentDetail);
		}

		let orderSummary = new OrderSummary();
		orderSummary.summary = order;

		let htmlPaymentCompletePage = createElement('div', {
			class:'shop-payment-complete-page',
			childs: [
				createElement('div', {
					class:'shop-payment-complete-header',
					i18n:'#payment_complete',
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
				}),*/
				createElement('label-property', {
					label: '#order_id',
					property: order.id,
				}),
				transactionSummary(paymentDetail),
				createElement('div', {
					class:'shop-payment-complete-content',
					childs: [
						createElement('div', {
							class:'shop-payment-complete-cart',
							child: orderSummary.html,
						}),

					]
				}),
			]
		});


		this.htmlContent.append(htmlPaymentCompletePage);
	}

	async #initCheckout() {
		const { requestId, response } = await fetchApi({ action: 'init-checkout', version: 1 });
		if (response?.success) {
			const order = new Order();
			order.fromJSON(response.result.order);
			this.#appContent.setOrder(order);

			//this.#orderView.model = order;
			this.#order = order;
			//this.#orderView.step = 'init';

			this.#orderSummary.summary = order;
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
			this.#orderSummary.summary = order;
			if (this.#pageSubType == PAGE_SUBTYPE_CHECKOUT_PAYMENT) {
				//this.#orderView.initPaypal(this.#orderId);
			}
		}
	}

	#initAddress() {
		if (!this.#order) {
			this.#navigateTo('/@checkout');
			return;
		}
		//this.#orderView.step = 'address';

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



		//this.#orderView.step = 'shipping';

		this.#displayCheckout();
	}

	async #sendShippingAddress() {
		const { requestId, response } = await fetchApi({
			action: 'set-shipping-address',
			version: 1,
			params: {
				shipping_address: this.#order.shippingAddress,
				//orderId: this.#orderId
			},
		});


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



		if (response?.success) {
			this.#order.fromJSON(response.result.order);
			this.#appContent.setOrder(this.#order);
			return true;
		}
	}

	async #sendShippingMethod() {
		const { requestId, response } = await fetchApi({ action: 'setshippingmethod', version: 1, method: this.#order.shippingMethod, orderId: this.#orderId });
		if (response?.success) {
			this.#order.fromJSON(response.result.order);
			this.#appContent.setOrder(this.#order);
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
			Controller.dispatchEvent(new CustomEvent('addnotification', {detail: {type: 'error', content: createElement('span', {i18n:'#error_while_creating_order'})}}));
			return;
		}

		let result = await this.#createPrintfulOrder();
		if (TESTING) {
			console.log(result);
		}
		if (result) {
			//this.#orderView.step = 'payment';
			this.#displayCheckout();
		}
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
		this.#broadcastChannel.postMessage({action: 'reloadcart'});
	}

	async #createPrintfulOrder() {
		this.#printfulOrder = null;

		const { requestId, response } = await fetchApi({ action: 'createprintfulorder', version: 1, orderID: this.#orderId });
		if (response?.success) {
			this.#printfulOrder = response.result.order;
			return this.#printfulOrder;
		} else {
			Controller.dispatchEvent(new CustomEvent('addnotification', {detail: {type: 'error', content: createElement('span', {
					'i18n-json': {
						innerHTML: '#error_while_creating_order',
					},
					'i18n-values': {
						requestId: requestId,
					},
				}),
			}}));
		}
		return false;
	}

	#getShopProductElement() {
		if (!this.#htmlShopProduct) {
			this.#htmlShopProduct = createElement('shop-product');
		}
		return this.#htmlShopProduct;
	}

	async #initProductPage(productId) {
		const product = await getShopProduct(productId);
		if (product) {
			this.#appContent.setProduct(product);
		}
	}

	async #initOrderPage(orderId) {
		const { requestId, response } = await fetchApi({
			action: 'getorder',
			version: 1,
			orderId: orderId,
		});
		if (response && response.success) {
			this.#viewOrderPage(response.result);
		} else {
			Controller.dispatchEvent(new CustomEvent('addnotification', {detail: {type: 'error', content: createElement('span', {i18n:'#failed_to_get_order_details'})}}));
		}
	}

	async #viewOrderPage(order) {
		createElement('div', {
			class:'shop-order-page',
			parent: this.htmlContent,
			child: orderSummary(order),
		});
	}

	async #displayProducts() {
		const shopProducts = await this.#refreshProducts();

		this.#appContent.setProducts(shopProducts);
	}

	async #refreshProducts() {
		const { requestId, response } = await fetchApi({
			action: 'get-products',
			version: 1,
		});

		console.log(response);

		if (response?.success) {
			console.log(response);
			const shopProducts = [];
			for (const productJSON of response.result) {
				const shopProduct = new ShopProduct();
				shopProduct.fromJSON(productJSON);
				shopProducts.push(shopProduct);
			}
			return shopProducts;
		} else {
			//Controller.dispatchEvent(new CustomEvent('addnotification', {detail: {type: 'error', content: createElement('span', {i18n:'#error_while_sending_message'})}}));
			//this.#htmlContactButton.disabled = false;
		}
	}

	async #displayFavorites() {
		this.#appContent.setFavorites(this.#favorites);

		return;
		let htmlFavoritesPage = createElement('div', {class:'shop-favorites-page', parent: this.htmlContent});
		let htmlFavoriteDetail = createElement('div', {class:'shop-favorites-detail shop-block', parent: htmlFavoritesPage});
		this.#htmlFavoriteList = createElement('div', {class:'shop-favorites-list', parent: htmlFavoriteDetail});
		await this.#refreshFavorites();
	}

	#initPage() {
		this.#htmlElement = createElement('div', {
			parent: document.body,
			attachShadow: { mode: 'closed' },
			adoptStyle: applicationCSS,
			childs:[
				this.#appToolbar.htmlElement,
				this.#appContent.htmlElement,
				this.#appFooter.htmlElement,
			]
		});

		//this.this.#htmlElement();
		//this.initContent();
		//this.initFooter();
	}

	async #initSession() {
		//let response = await fetch('/getcurrency');
		//let json = await response.json();
		const result = await ServerAPI.getCurrency();
		//if (json && json.success) {
		this.#setCurrency(result);
		//}
	}

	#setCurrency(currency) {
		this.#currency = currency;
		this.#appToolbar.setCurrency(currency);

		//this.#htmlCurrency.innerHTML = `${I18n.getString('#currency')} ${currency}`;
	}

	async #setServerCurrency(currency) {
		let result = await this.#post('/setcurrency', {currency:currency});
		if (result && result.success) {
			if (this.#pageType == PAGE_TYPE_CHECKOUT) {
				this.#navigateTo('/@checkout');
			}
		}
	}

	async #get(url) {
		let response = await fetch(url);
		return await response.json();
	}

	async #post(url, body) {
		let response = await fetch(url, {method:'POST', body:JSON.stringify(body)});
		return await response.json();
	}

	#navigateTo(url, replaceSate  = false) {
		history[replaceSate ? 'replaceState' : 'pushState']({}, undefined, url);
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

	#pushState(url) {
		history.pushState({}, undefined, url);
	}

	#replaceState(url) {
		history.replaceState(this.#getHistoryState(), undefined, url);
	}

	#historyStateChanged() {
		history.replaceState(this.#getHistoryState(), '');
	}

	#getHistoryState() {
		return {
			columnCartVisible: this.#htmlColumnCartVisible,
		};
	}

	#restoreHistoryState({ columnCartVisible = false } = {}) {
		this.#htmlColumnCartVisible = columnCartVisible;
	}

	#onPaymentComplete(order) {
		this.#order.fromJSON(order);
		console.log(this.#order);
		this.#paymentCompleteDetails = {order: this.#order};
		this.#order = null;
		//this.#orderView.model = null;
		this.#orderSummary.summary = null;
		this.#loadCart();

		this.#navigateTo(`/@order/${order.id}`);
	}

	initHeader() {
		let htmlHeader = createElement('header', {class:'shop-header'});
		this.htmlHeader = htmlHeader;
		this.#htmlElement.append(htmlHeader);

		createElement('div', {
			class: 'shop-header-font-size',
			parent: htmlHeader,
			childs: [
				createElement('div', {
					class: 'smaller',
					innerHTML: textDecreaseSVG,
					events: {
						click: () => {
							this.#fontSize /= 1.1;
							document.documentElement.style.setProperty('--font-size', `${this.#fontSize}px`);
						}
					}
				}),
				createElement('div', {
					class: 'larger',
					innerHTML: textIncreaseSVG,
					events: {
						click: () => {
							this.#fontSize *= 1.1;
							document.documentElement.style.setProperty('--font-size', `${this.#fontSize}px`);
						}
					}
				}),
			]
		});

		let htmlCurrencyContainer = createElement('div', {class:'shop-header-currency'});
		this.#htmlCurrency = createElement('div');
		let htmlCurrencySelector = createElement('ul', {class:'shop-header-currency-list',hidden:1});
		document.addEventListener('click', (event) => {
			//event.target == this.#htmlCurrency ? show(htmlCurrencySelector) : hide(htmlCurrencySelector);
		});
		htmlCurrencyContainer.append(this.#htmlCurrency, htmlCurrencySelector);
		for(let currency of ALLOWED_CURRENCIES) {
			let htmlCurrencyItem = createElement('li', {class:'shop-header-currency-item'});
			htmlCurrencyItem.innerHTML = currency;
			htmlCurrencyItem.addEventListener('click', () => {this.#setCurrency(currency);this.#setServerCurrency(currency);});

			htmlCurrencySelector.append(htmlCurrencyItem);
		}

		this.htmlHeaderProducts = createElement('div', {
			class: 'shop-header-products',
			parent: htmlHeader,
			i18n: '#products',
			events: {
				click: () => this.#navigateTo('/@products'),
				mouseup: (event) => {
					if (event.button == 1) {
						open('@products', '_blank');
					}
				},
			}
		});

		createElement('div', {
			class: 'shop-header-favorites',
			parent: htmlHeader,
			childs: [
				createElement('div', {
					class: 'favorites-icon',
					innerHTML: bookmarksPlainSVG,
				}),
				this.#htmlHeaderFavorites = createElement('div', {
					class: 'favorites-count',
				}),
			],
			events: {
				click: () => this.#navigateTo('/@favorites'),
				mouseup: (event) => {
					if (event.button == 1) {
						open('@favorites', '_blank');
					}
				},
			}
		});

		createElement('div', {
			class:'shop-header-cart',
			parent: htmlHeader,
			childs: [
				createElement('div', {
					class: 'cart-icon',
					innerHTML: shoppingCartSVG,
				}),
				this.#htmlCart = createElement('div', {
					class: 'cart-count',
				}),
			],
			events: {
				click: () => this.#navigateTo('/@cart'),
				mouseup: (event) => {
					if (event.button == 1) {
						open('@cart', '_blank');
					}
				},
			}
		});

		let htmlTheme = createElement('input', {class:'shop-header-theme', type:'checkbox'});
		htmlTheme.addEventListener('input', () => this.theme = htmlTheme.checked ? 'dark' : 'light');

		htmlHeader.append(htmlCurrencyContainer, htmlTheme);
	}

	initContent() {
		let htmlContent = createElement('section', {class:'shop-content'});
		this.htmlContent = htmlContent;
		this.#htmlElement.append(htmlContent);
		this.#htmlProductsPage = createElement('div', {class:'shop-product-list'});

		this.#htmlCartList = createElement('cart-products');
	}

	initFooter() {
		this.#htmlFooter = createElement('footer', {
			class:'shop-footer',
			parent: this.#htmlElement,
			childs: [
				createElement('span', {
					i18n: '#contact',
					events: {
						click: () => this.#navigateTo('/@contact'),
						mouseup: (event) => {
							if (event.button == 1) {
								open('@contact', '_blank');
							}
						},
					}
				}),
				createElement('span', {
					i18n: '#privacy_policy',
					events: {
						click: () => this.#navigateTo('/@privacy'),
						mouseup: (event) => {
							if (event.button == 1) {
								open('@privacy', '_blank');
							}
						},
					}
				}),
				createElement('span', {
					i18n: '#cookies',
					events: {
						click: () => this.#navigateTo('/@cookies'),
						mouseup: (event) => {
							if (event.button == 1) {
								open('@cookies', '_blank');
							}
						},
					}
				}),
			]
		});
	}

	set theme(theme) {
		document.documentElement.classList.remove('light');
		document.documentElement.classList.remove('dark');
		//document.documentElement.classList.add(theme);
	}

	async #processMessage(event) {
		switch (event.data.action) {
			case 'cartchanged':
				this.#cart.fromJSON(event.data.cart);
				const showColumnCart = this.#cart.totalQuantity > 0;
				this.#htmlColumnCartVisible = showColumnCart;
				this.#htmlColumnCart?.display(showColumnCart);
				this.#historyStateChanged();
				Controller.dispatchEvent(new CustomEvent(EVENT_REFRESH_CART, { detail: this.#cart }));
				Controller.dispatchEvent(new CustomEvent(EVENT_CART_COUNT, { detail: this.#cart.totalQuantity }));
				break;
			case 'cartloaded':
				Controller.dispatchEvent(new CustomEvent(EVENT_CART_COUNT, { detail: this.#cart.totalQuantity }));
				break;
			case 'reloadcart':
				this.#loadCart();
				break;
			case 'favoriteschanged':
				this.#favorites = event.data.favorites;
				await this.#refreshFavorites();
				this.#countFavorites();
				break;
		}
	}

	async #loadCart() {
		const { requestId, response } = await fetchApi({
			action: 'get-cart',
			version: 1,
		});
		if (TESTING) {
			console.log(response);
		}

		this.#cart.fromJSON(response?.result?.cart);

		this.#refreshCart();

		this.#broadcastChannel.postMessage({ action: 'cartloaded', cart: this.#cart.toJSON() });
	}
}
new Application();
