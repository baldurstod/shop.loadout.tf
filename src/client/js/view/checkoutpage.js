import { createElement, hide, show } from 'harmony-ui';
import { CheckoutAddresses } from './checkoutaddresses.js';
import { ShippingMethodSelector } from './shippingmethodselector.js';
export { Address } from './components/address.js';

import checkoutPageCSS from '../../css/checkoutpage.css';
import { PAGE_SUBTYPE_CHECKOUT_ADDRESS, PAGE_SUBTYPE_CHECKOUT_INIT, PAGE_SUBTYPE_CHECKOUT_SHIPPING } from '../constants.js';

export class CheckoutPage {
	#htmlElement;
	#checkoutAddress = new CheckoutAddresses();
	#shippingMethodSelector = new ShippingMethodSelector();

	constructor() {
		this.#initHTML();
	}

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: checkoutPageCSS,
			childs: [
				this.#checkoutAddress.htmlElement,
				this.#shippingMethodSelector.htmlElement,
				//this.#htmlShippingAddress = createElement('shop-address'),
				//this.#htmlBillingAddress = createElement('shop-address'),
			],
		});
		return this.#htmlElement;
	}

	setCheckoutStage(pageSubType) {
		hide(this.#checkoutAddress.htmlElement);
		hide(this.#shippingMethodSelector.htmlElement);
		switch (pageSubType) {
			case PAGE_SUBTYPE_CHECKOUT_INIT:
				break;
			case PAGE_SUBTYPE_CHECKOUT_ADDRESS:
				show(this.#checkoutAddress.htmlElement);
				break;
			case PAGE_SUBTYPE_CHECKOUT_SHIPPING:
				show(this.#shippingMethodSelector.htmlElement);
				break;
			default:
				throw `Unknown page type ${pageSubType}`;
				break;
		}
	}

	setOrder(order) {
		this.#checkoutAddress.setOrder(order);
		this.#shippingMethodSelector.setOrder(order);
		//this.#htmlShippingAddress.setAddress(order.shippingAddress);
		//this.#htmlBillingAddress.setAddress(order.billingAddress);
	}

	setCountries(countries) {
		this.#checkoutAddress.setCountries(countries);
		//this.#htmlShippingAddress.setCountries(countries);
		//this.#htmlBillingAddress.setCountries(countries);
	}

	get htmlElement() {
		return this.#htmlElement;
	}
}
