import { I18n, createElement, createShadowRoot } from 'harmony-ui';
import { getCartTotalPriceFormatted } from '../carttotalprice';
import { Controller } from '../controller';
import { EVENT_NAVIGATE_TO } from '../controllerevents';
import commonCSS from '../../css/common.css';
import cartPageCSS from '../../css/cartpage.css';
import { defineCartProducts, HTMLCartProductsElement } from './components/cartproducts';
import { Cart } from '../model/cart';

export class CartPage {
	#shadowRoot?: ShadowRoot;
	#htmlDetail?: HTMLElement;
	#htmlCartList?: HTMLCartProductsElement;
	#htmlSubtotalLine?: HTMLElement;
	#htmlSubtotalLabel?: HTMLElement;
	#htmlSubtotal?: HTMLElement;
	#htmlCheckout?: HTMLElement;
	#htmlCheckoutSubtotalLabel?: HTMLElement;
	#htmlCheckoutSubtotal?: HTMLElement;
	#htmlCheckoutButton?: HTMLButtonElement;

	#initHTML() {
		defineCartProducts();
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyles: [commonCSS, cartPageCSS],
			childs: [
				this.#htmlDetail = createElement('div', {
					class: 'detail',
					childs: [
						createElement('div', {
							i18n: '#shoppingcart',
							class: 'header'
						}),
						this.#htmlCartList = createElement('cart-products') as HTMLCartProductsElement,
						this.#htmlSubtotalLine = createElement('div', {
							class: 'subtotal shop-cart-line',
							childs: [
								this.#htmlSubtotalLabel = createElement('span', {
									class: 'label',
									'i18n-json': {
										innerHTML: '#subtotal_count',
									},
									'i18n-values': {
										count: 0,
									},
								}),
								this.#htmlSubtotal = createElement('span', { class: 'price' }),
							]
						}),
					],
				}),
				this.#htmlCheckout = createElement('div', {
					class: 'checkout',
					childs: [
						this.#htmlCheckoutSubtotalLabel = createElement('span', {
							class: 'label',
							'i18n-json': {
								innerHTML: '#subtotal_count',
							},
							'i18n-values': {
								count: 0,
							},
						}),
						this.#htmlCheckoutSubtotal = createElement('span', { class: 'price' }),
						this.#htmlCheckoutButton = createElement('button', {
							i18n: '#checkout',
							events: {
								click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout' } })),
							}
						}) as HTMLButtonElement,
					],
				}),
			],
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	get htmlElement() {
		throw 'use getHTML';
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}

	async setCart(cart: Cart) {
		if (!this.#shadowRoot) {
			this.#initHTML();
		}

		this.#htmlCartList!.setCart(cart);
		this.#htmlSubtotalLabel!.setAttribute('data-i18n-values', JSON.stringify({ count: cart.totalQuantity, }));
		this.#htmlCheckoutSubtotalLabel!.setAttribute('data-i18n-values', JSON.stringify({ count: cart.totalQuantity, }));
		this.#htmlCheckoutSubtotal!.innerText = this.#htmlSubtotal!.innerText = await getCartTotalPriceFormatted(cart);
	}
}
