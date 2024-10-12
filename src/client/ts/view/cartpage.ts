import { I18n, createElement } from 'harmony-ui';
export { CartProducts } from './components/cartproducts.js';
import { getCartTotalPriceFormatted } from '../carttotalprice.js';
import { Controller } from '../controller.js';
import { EVENT_NAVIGATE_TO } from '../controllerevents.js';

import commonCSS from '../../css/common.css';
import cartPageCSS from '../../css/cartpage.css';

export class CartPage {
	#htmlElement;
	#htmlDetail;
	#htmlCartList;
	#htmlSubtotalLine;
	#htmlSubtotalLabel;
	#htmlSubtotal;
	#htmlCheckout;
	#htmlCheckoutSubtotalLabel;
	#htmlCheckoutSubtotal;
	#htmlCheckoutButton;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [ commonCSS, cartPageCSS ],
			childs: [
				this.#htmlDetail = createElement('div', {
					class: 'detail',
					childs: [
						createElement('div', {
							i18n: '#shoppingcart',
							class: 'header'
						}),
						this.#htmlCartList = createElement('cart-products'),
						this.#htmlSubtotalLine = createElement('div', {
							class: 'subtotal shop-cart-line',
							childs:[
								this.#htmlSubtotalLabel = createElement('span', {
									class: 'label',
									'i18n-json': {
										innerHTML: '#subtotal_count',
									},
									'i18n-values': {
										count: 0,
									},
								}),
								this.#htmlSubtotal = createElement('span', { class:'price' }),
							]
						}),
					],
				}),
				this.#htmlCheckout = createElement('div', {
					class:'checkout',
					childs:[
						this.#htmlCheckoutSubtotalLabel = createElement('span', {
							class: 'label',
							'i18n-json': {
								innerHTML: '#subtotal_count',
							},
							'i18n-values': {
								count: 0,
							},
						}),
						this.#htmlCheckoutSubtotal = createElement('span', { class:'price' }),
						this.#htmlCheckoutButton = createElement('button', {
							i18n: '#checkout',
							events: {
								click: () => Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout' } })),
							}
						}),
					],
				}),
			],
		});

		I18n.observeElement(this.#htmlElement);

		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}

	async setCart(cart) {
		this.#htmlCartList.setCart(cart);

		this.#htmlSubtotalLabel.setAttribute('data-i18n-values', JSON.stringify({ count: cart.totalQuantity, }));
		this.#htmlCheckoutSubtotalLabel.setAttribute('data-i18n-values', JSON.stringify({ count: cart.totalQuantity, }));
		this.#htmlCheckoutSubtotal.innerHTML = this.#htmlSubtotal.innerHTML = await getCartTotalPriceFormatted(cart);
	}
}
