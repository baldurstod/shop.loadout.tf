import { createElement } from 'harmony-ui';
import { Cart } from '../../model/cart';
import { defineCartItem, HTMLCartItemElement } from './cartitem';

export class HTMLCartProductsElement extends HTMLElement {

	#refreshHTML(cart: Cart): void {
		this.innerText = '';
		defineCartItem();

		if (cart.totalQuantity > 0) {
			for (const [productID, quantity] of cart.items) {
				this.append(createElement('cart-item', {
					elementCreated: (element: Element) => { (element as HTMLCartItemElement).setItem(productID, quantity, cart.currency) },
				}));
			}
		} else {
			this.append(createElement('div', {
				i18n: '#empty_cart',
			}));
		}
	}

	setCart(cart: Cart): void {
		this.#refreshHTML(cart);
	}
}

let definedCartProducts = false;
export function defineCartProducts(): void {
	if (window.customElements && !definedCartProducts) {
		customElements.define('cart-products', HTMLCartProductsElement);
		definedCartProducts = true;
	}
}
