import { createElement, I18n } from 'harmony-ui';

export class CartProducts extends HTMLElement {

	#refreshHTML(cart) {
		this.innerHTML = '';

		if (cart.totalQuantity > 0) {
			for (let [productID, quantity] of cart.items) {
				this.append(createElement('cart-item', {
					elementCreated: element => element.setItem(productID, quantity, cart.currency),
				}));
			}
		} else {
			this.append(createElement('div', {
				i18n: '#empty_cart',
			}));
		}
	}

	setCart(cart) {
		this.#refreshHTML(cart);
	}
}

if (window.customElements) {
	customElements.define('cart-products', CartProducts);
}
