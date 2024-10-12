import { updateElement, createElement, shadowRootStyle, I18n } from 'harmony-ui';
import { Controller } from '../../controller.js';
import { formatPrice } from '../../utils.js';
import { getShopProduct } from '../../shopproducts.js';
import { MAX_PRODUCT_QTY } from '../../constants.js';
import { EVENT_NAVIGATE_TO } from '../../controllerevents.js';

import cartItemCSS from '../../../css/cartitem.css';
import commonCSS from '../../../css/common.css';

export class CartItemElement extends HTMLElement {
	#shadowRoot;
	#htmlProductName;
	#htmlProductThumb;
	#htmlProductPrice;
	#htmlProductQuantity;
	#productID;
	#quantity;
	#currency;
	#product;
	constructor() {
		super();
		this.#initHTML();
	}

	#initHTML() {
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
		shadowRootStyle(this.#shadowRoot, commonCSS);
		shadowRootStyle(this.#shadowRoot, cartItemCSS);

		this.#htmlProductName = createElement('div', {
			class: 'name',
			parent: this.#shadowRoot,
			events: {
				aclick: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: product.shopUrl } })),
				click: () => console.info(this),
				mouseup: (event) => {
					if (event.button == 1) {
						open(product.shopUrl, '_blank');
					}
				},
			}
		});

		this.#htmlProductThumb = createElement('img', {
			class: 'thumb',
			parent: this.#shadowRoot,
			events: {
				click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: product.shopUrl } })),
				mouseup: (event) => {
					if (event.button == 1) {
						open(product.shopUrl, '_blank');
					}
				},
			}
		});

		createElement('td', {
			class:'infos',
			parent: this.#shadowRoot,
			childs: [
				this.#htmlProductQuantity = createElement('input', {
					class: 'quantity',
					type: 'number',
					min: 1,
					max: MAX_PRODUCT_QTY,
					events: {
						input: (event) => {
							let q = Number(event.target.value);
							if (!Number.isNaN(q) && q > 0) {
								Controller.dispatchEvent(new CustomEvent('setquantity', { detail: { id: this.#productID, quantity: q } }));
							}
						}
					}
				}),
				createElement('button', {
					class: 'remove',
					innerHTML: '🗑️',
					events: {
						click: () => {
							Controller.dispatchEvent(new CustomEvent('setquantity', { detail: { id: this.#productID, quantity: 0 } }));
						}
					}
				}),
			]
		});

		this.#htmlProductPrice = createElement('td', {
			class: 'price',
			parent: this.#shadowRoot,
		});

/*
		this.#htmlCartCheckout = createElement('div', {
			parent: this.#shadowRoot,
			class:'checkout',
			hidden: true,
			childs: [
				this.#htmlSubtotalLabel = createElement('span', {class:'shop-cart-checkout-subtotal-label'}),
				this.#htmlSubtotal = createElement('span', {class:'shop-price'}),
				this.#htmlGotoCart = createElement('div', {
					class:'goto-cart',
					i18n: '#go_to_cart',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, {detail:{url:'/@cart'}})),
					}
				}),
				this.#htmlCheckout = createElement('button', {
					class:'shop-cart-checkout-button',
					i18n: '#checkout',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout' } })),
					}
				}),
				this.#htmlItemList = createElement('div', {class:'item-list'}),
			],
		});*/
	}

	#refreshHTML(/*product, currency*/) {
		const product = this.#product;
		const currency = this.#currency;

		if (!product) {
			return;
		}
		updateElement(this, {
			class: 'shop-cart-line',
			events: {
				click: event => {
					if (event.target == event.currentTarget) {
						Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, {detail: {url: product.shopUrl}}))
					}
				},
				mouseup: event => {
					if (event.target == event.currentTarget && event.button == 1) {
						open(product.shopUrl, '_blank');
					}
				},
			}
		});


		this.#htmlProductName.innerText = product.name;
		this.#htmlProductThumb.src = product.thumbnailUrl;
		this.#htmlProductPrice.innerText = formatPrice(product.retailPrice, currency);
		this.#htmlProductQuantity.value = this.#quantity;

		//htmlElement.append(htmlProductThumb, htmlProductInfo, htmlProductPrice);
	}

	async setItem(productID, quantity, currency) {
		this.#productID = productID;
		this.#quantity = quantity;
		this.#currency = currency;
		this.#product = await getShopProduct(productID);
		this.#refreshHTML(/*product, currency*/);
	}
}

if (window.customElements) {
	customElements.define('cart-item', CartItemElement);
}
