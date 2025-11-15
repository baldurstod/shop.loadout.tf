import { createElement, display, hide, I18n, shadowRootStyle, show } from 'harmony-ui';
import columnCartCSS from '../../../css/columncart.css';
import commonCSS from '../../../css/common.css';
import { getCartTotalPriceFormatted } from '../../carttotalprice';
import { Controller, ControllerEvent, NavigateToDetail } from '../../controller';
import { Cart } from '../../model/cart';
import { defineCartItem, HTMLCartItemElement } from './cartitem';

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
		Controller.addEventListener(ControllerEvent.RefreshCart, (event: Event) => { this.#refreshHTML((event as CustomEvent<Cart>).detail) });
	}

	#initHTML(): void {
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
						click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@cart' } }),
					}
				}),
				this.#htmlCheckout = createElement('button', {
					class: 'shop-cart-checkout-button',
					i18n: '#checkout',
					events: {
						click: () => Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@checkout' } }),
					}
				}) as HTMLButtonElement,
				this.#htmlItemList = createElement('div', { class: 'item-list' }),
			],
		});
	}

	async #refreshHTML(cart: Cart): Promise<void> {
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
		for (const [productID, quantity] of cart.items) {
			//this.#htmlItemList.append(product.toHTML(cart.currency));
			//this.#htmlItemList.append(htmlCartProduct = createElement('cart-product'));
			//htmlCartProduct.setProduct(product, cart.currency);

			createElement('cart-item', {
				parent: this.#htmlItemList,
				elementCreated: (element: Element) => { (element as HTMLCartItemElement).setItem(productID, quantity, cart.currency) },
			})
		}
	}

	/*set label(label) {
		this.#htmlLabel.setAttribute('data-i18n', label);
	}

	set property(property) {
		this.#htmlProperty.innerText = property;
	}*/

	display(visible: boolean): void {
		display(this, visible);
	}
}

let definedColumnCart = false;
export function defineColumnCart(): void {
	if (window.customElements && !definedColumnCart) {
		customElements.define('column-cart', HTMLColumnCartElement);
		definedColumnCart = true;
	}
}
