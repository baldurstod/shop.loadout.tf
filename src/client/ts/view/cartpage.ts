import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import cartPageCSS from '../../css/cartpage.css';
import commonCSS from '../../css/common.css';
import { getCartTotalPriceFormatted } from '../carttotalprice';
import { Controller } from '../controller';
import { EVENT_NAVIGATE_TO } from '../controllerevents';
import { Cart } from '../model/cart';
import { defineCartProducts, HTMLCartProductsElement } from './components/cartproducts';
import { ShopElement } from './shopelement';

export class CartPage extends ShopElement {
	#htmlDetail?: HTMLElement;
	#htmlCartList?: HTMLCartProductsElement;
	#htmlSubtotalLine?: HTMLElement;
	#htmlSubtotalLabel?: HTMLElement;
	#htmlSubtotal?: HTMLElement;
	#htmlCheckout?: HTMLElement;
	#htmlCheckoutSubtotalLabel?: HTMLElement;
	#htmlCheckoutSubtotal?: HTMLElement;
	#htmlCheckoutButton?: HTMLButtonElement;

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}
		defineCartProducts();
		this.shadowRoot = createShadowRoot('section', {
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
									i18n: {
										innerText: '#subtotal_count',
										values: {
											count: 0,
										},
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
							i18n: {
								innerText: '#subtotal_count',
								values: {
									count: 0,
								},
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
		I18n.observeElement(this.shadowRoot);
	}

	async setCart(cart: Cart): Promise<void> {
		this.initHTML();

		this.#htmlCartList!.setCart(cart);
		I18n.setValue(this.#htmlSubtotalLabel, 'count', cart.totalQuantity);
		I18n.setValue(this.#htmlCheckoutSubtotalLabel, 'count', cart.totalQuantity);
		this.#htmlCheckoutSubtotal!.innerText = this.#htmlSubtotal!.innerText = await getCartTotalPriceFormatted(cart);
	}
}
