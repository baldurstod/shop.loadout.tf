import { createElement } from 'harmony-ui';
import { defineCartItem, HTMLCartItemElement } from './cartitem';
import { Cart } from '../../model/cart';

export class HTMLCartProductsElement extends HTMLElement {

	#refreshHTML(cart: Cart) {
		this.innerText = '';
		defineCartItem();

		if (cart.totalQuantity > 0) {
			for (let [productID, quantity] of cart.items) {
				this.append(createElement('cart-item', {
					elementCreated: (element: HTMLElement) => (element as HTMLCartItemElement).setItem(productID, quantity, cart.currency),
				}));
			}
		} else {
			this.append(createElement('div', {
				i18n: '#empty_cart',
			}));
		}
	}

	setCart(cart: Cart) {
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
