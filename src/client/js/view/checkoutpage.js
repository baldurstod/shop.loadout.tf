import { createElement } from 'harmony-ui';
export { Address } from './components/address.js';

import checkoutPageCSS from '../../css/checkoutpage.css';

export class CheckoutPage {
	#htmlElement;
	#htmlShippingAddress;
	#htmlBillingAddress;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: checkoutPageCSS,
			childs: [
				this.#htmlShippingAddress = createElement('shop-address'),
				this.#htmlBillingAddress = createElement('shop-address'),
			],
		});
		return this.#htmlElement;
	}

	setSubType() {

	}

	setOrder(order) {
		this.#htmlShippingAddress.setAddress(order.shippingAddress);
		this.#htmlBillingAddress.setAddress(order.billingAddress);
	}

	setCountries(countries) {
		this.#htmlShippingAddress.setCountries(countries);
		this.#htmlBillingAddress.setCountries(countries);
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}
}
