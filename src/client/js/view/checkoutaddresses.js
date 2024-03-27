import { I18n, createElement, display } from 'harmony-ui';
import 'harmony-ui/dist/define/harmony-switch.js';
export { Address } from './components/address.js';
import { Controller } from '../controller.js';
import { EVENT_NAVIGATE_TO } from '../controllerevents.js';

import checkoutAddressesCSS from '../../css/checkoutaddresses.css';
import commonCSS from '../../css/common.css';

export class CheckoutAddresses {
	#htmlElement;
	#htmlShippingAddress;
	#htmlSameBillingAddress;
	#htmlBillingAddress;
	#htmlContinue;
	#order;

	constructor() {
		this.#initHTML();
	}

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [ checkoutAddressesCSS, commonCSS ],
			childs: [
				this.#htmlShippingAddress = createElement('shop-address', {
					elementCreated: element => element.setAddressType('#shipping_address'),
				}),
				this.#htmlSameBillingAddress = createElement('harmony-switch', {
					i18n: '#same_billing_address',
					events: {
						change: event => this.#changeSameBillingAddress(event.target.state),
					},
				}),
				this.#htmlBillingAddress = createElement('shop-address', {
					elementCreated: element => element.setAddressType('#billing_address'),
				}),
				this.#htmlContinue = createElement('button', {
					i18n: '#continue_to_shipping',
					events: {
						click: () => this.#continueCheckout(),
					},
				}),
			],
		});
		I18n.observeElement(this.#htmlElement);
	}

	#refresh() {
		if (!this.#order) {
			return;
		}

		const sameBillingAddress = this.#order?.getSameBillingAddress();
		this.#htmlSameBillingAddress.state = sameBillingAddress;

		display(this.#htmlBillingAddress, !sameBillingAddress);
	}

	setOrder(order) {
		this.#order = order;
		this.#htmlShippingAddress.setAddress(order.shippingAddress);
		this.#htmlBillingAddress.setAddress(order.billingAddress);
		this.#refresh();
	}

	setCountries(countries) {
		this.#htmlShippingAddress.setCountries(countries);
		this.#htmlBillingAddress.setCountries(countries);
	}

	#changeSameBillingAddress(sameBillingAddress) {
		this.#order?.setSameBillingAddress(sameBillingAddress);
		this.#refresh();
	}

	#continueCheckout() {
		//TODO: check values
		Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout#shipping' } }));
	}

	get htmlElement() {
		return this.#htmlElement;
	}
}
