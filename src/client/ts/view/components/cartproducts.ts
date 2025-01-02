import { createElement, I18n } from 'harmony-ui';
import { defineCartItem } from './cartitem';

export class HTMLCartProductsElement extends HTMLElement {

	#refreshHTML(cart) {
		this.innerHTML = '';
		defineCartItem();

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

let definedCartProducts = false;
export function defineCartProducts() {
	if (window.customElements && !definedCartProducts) {
		customElements.define('cart-products', HTMLCartProductsElement);
		definedCartProducts = true;
	}
}
