import { I18n, createElement, hide, shadowRootStyle, show } from 'harmony-ui';
import { Controller } from '../../controller';
import { getCartTotalPriceFormatted } from '../../carttotalprice';
import { EVENT_NAVIGATE_TO, EVENT_REFRESH_CART } from '../../controllerevents';
import columnCartCSS from '../../../css/columncart.css';
import commonCSS from '../../../css/common.css';
import { defineCartItem } from './cartitem';

export class HTMLColumnCartElement extends HTMLElement {
	#shadowRoot;
	#htmlCartCheckout;
	#htmlSubtotalLabel;
	#htmlSubtotal;
	#htmlGotoCart;
	#htmlItemList;
	#htmlCheckout;
	constructor() {
		super();
		this.#initHTML();
		Controller.addEventListener(EVENT_REFRESH_CART, (event: CustomEvent) => this.#refreshHTML(event.detail));
	}

	#initHTML() {
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
		shadowRootStyle(this.#shadowRoot, commonCSS);
		shadowRootStyle(this.#shadowRoot, columnCartCSS);

		this.#htmlCartCheckout = createElement('div', {
			parent: this.#shadowRoot,
			class: 'checkout',
			hidden: true,
			childs: [
				this.#htmlSubtotalLabel = createElement('span', { class: 'shop-cart-checkout-subtotal-label' }),
				this.#htmlSubtotal = createElement('span', { class: 'price' }),
				this.#htmlGotoCart = createElement('div', {
					class: 'goto-cart',
					i18n: '#go_to_cart',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@cart' } })),
					}
				}),
				this.#htmlCheckout = createElement('button', {
					class: 'shop-cart-checkout-button',
					i18n: '#checkout',
					events: {
						click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout' } })),
					}
				}),
				this.#htmlItemList = createElement('div', { class: 'item-list' }),
			],
		});
	}

	async #refreshHTML(cart) {
		if (cart.totalQuantity > 0) {
			show(this.#htmlCartCheckout);
		} else {
			hide(this.#htmlCartCheckout);
			return;
		}
		this.#htmlSubtotalLabel.innerHTML = `${I18n.getString('#subtotal')} (${cart.totalQuantity} ${cart.totalQuantity > 1 ? I18n.getString('#items') : I18n.getString('#item')}): `;
		this.#htmlSubtotal.innerHTML = await getCartTotalPriceFormatted(cart);//cart.totalPriceFormatted;

		this.#htmlItemList.innerHTML = '';
		defineCartItem();
		//let htmlCartProduct;
		for (let [productID, quantity] of cart.items) {
			//this.#htmlItemList.append(product.toHTML(cart.currency));
			//this.#htmlItemList.append(htmlCartProduct = createElement('cart-product'));
			//htmlCartProduct.setProduct(product, cart.currency);

			createElement('cart-item', {
				parent: this.#htmlItemList,
				elementCreated: element => element.setItem(productID, quantity, cart.currency),
			})
		}
	}

	/*set label(label) {
		this.#htmlLabel.setAttribute('data-i18n', label);
	}

	set property(property) {
		this.#htmlProperty.innerHTML = property;
	}*/

	display(display) {
		if (display) {
			show(this);
		} else {
			hide(this);
		}
	}
}

let definedColumnCart = false;
export function defineColumnCart() {
	if (window.customElements && !definedColumnCart) {
		customElements.define('column-cart', HTMLColumnCartElement);
		definedColumnCart = true;
	}
}
