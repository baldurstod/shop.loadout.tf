import { I18n, createElement, display, hide, shadowRootStyle, show } from 'harmony-ui';
import { Controller } from '../../controller';
import { getCartTotalPriceFormatted } from '../../carttotalprice';
import { EVENT_NAVIGATE_TO, EVENT_REFRESH_CART } from '../../controllerevents';
import columnCartCSS from '../../../css/columncart.css';
import commonCSS from '../../../css/common.css';
import { defineCartItem, HTMLCartItemElement } from './cartitem';
import { Cart } from '../../model/cart';

export class HTMLColumnCartElement extends HTMLElement {
	#shadowRoot!: ShadowRoot;
	#htmlCartCheckout!: HTMLElement;
	#htmlSubtotalLabel!: HTMLElement;
	#htmlSubtotal!: HTMLElement;
	#htmlGotoCart!: HTMLElement;
	#htmlItemList!: HTMLElement;
	#htmlCheckout!: HTMLButtonElement;
	constructor() {
		super();
		this.#initHTML();
		Controller.addEventListener(EVENT_REFRESH_CART, (event: Event) => this.#refreshHTML((event as CustomEvent).detail));
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
				}) as HTMLButtonElement,
				this.#htmlItemList = createElement('div', { class: 'item-list' }),
			],
		});
	}

	async #refreshHTML(cart: Cart) {
		if (cart.totalQuantity > 0) {
			show(this.#htmlCartCheckout);
		} else {
			hide(this.#htmlCartCheckout);
			return;
		}
		this.#htmlSubtotalLabel.innerText = `${I18n.getString('#subtotal')} (${cart.totalQuantity} ${cart.totalQuantity > 1 ? I18n.getString('#items') : I18n.getString('#item')}): `;
		this.#htmlSubtotal.innerText = await getCartTotalPriceFormatted(cart);//cart.totalPriceFormatted;

		this.#htmlItemList.innerText = '';
		defineCartItem();
		//let htmlCartProduct;
		for (let [productID, quantity] of cart.items) {
			//this.#htmlItemList.append(product.toHTML(cart.currency));
			//this.#htmlItemList.append(htmlCartProduct = createElement('cart-product'));
			//htmlCartProduct.setProduct(product, cart.currency);

			createElement('cart-item', {
				parent: this.#htmlItemList,
				elementCreated: (element: HTMLElement) => (element as HTMLCartItemElement).setItem(productID, quantity, cart.currency),
			})
		}
	}

	/*set label(label) {
		this.#htmlLabel.setAttribute('data-i18n', label);
	}

	set property(property) {
		this.#htmlProperty.innerText = property;
	}*/

	display(visible: boolean) {
		display(this, visible);
	}
}

let definedColumnCart = false;
export function defineColumnCart() {
	if (window.customElements && !definedColumnCart) {
		customElements.define('column-cart', HTMLColumnCartElement);
		definedColumnCart = true;
	}
}
