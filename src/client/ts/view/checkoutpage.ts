import { createElement, hide, show } from 'harmony-ui';
import { CheckoutAddresses } from './checkoutaddresses';
import { PaymentSelector } from './payment/paymentselector';
import { ShippingMethodSelector } from './shippingmethodselector';
import { PaypalPayment } from './payment/paypalpayment';

import checkoutPageCSS from '../../css/checkoutpage.css';
import { PAGE_SUBTYPE_CHECKOUT_ADDRESS, PAGE_SUBTYPE_CHECKOUT_INIT, PAGE_SUBTYPE_CHECKOUT_PAYMENT, PAGE_SUBTYPE_CHECKOUT_SHIPPING, PageSubType } from '../constants.js';

export class CheckoutPage {
	#htmlElement: HTMLElement;
	#checkoutAddress = new CheckoutAddresses();
	#shippingMethodSelector = new ShippingMethodSelector();
	#paymentSelector = new PaymentSelector();

	constructor() {
		this.#initHTML();

		this.#paymentSelector.addPaymentMethod(new PaypalPayment());
	}

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: checkoutPageCSS,
			childs: [
				this.#checkoutAddress.getHTML(),
				this.#shippingMethodSelector.htmlElement,
				this.#paymentSelector.htmlElement,
			],
		});
		return this.#htmlElement;
	}

	setCheckoutStage(pageSubType: PageSubType) {
		hide(this.#checkoutAddress.getHTML());
		hide(this.#shippingMethodSelector.htmlElement);
		hide(this.#paymentSelector.htmlElement);
		switch (pageSubType) {
			case PAGE_SUBTYPE_CHECKOUT_INIT:
				break;
			case PAGE_SUBTYPE_CHECKOUT_ADDRESS:
				show(this.#checkoutAddress.getHTML());
				break;
			case PAGE_SUBTYPE_CHECKOUT_SHIPPING:
				show(this.#shippingMethodSelector.htmlElement);
				break;
			case PAGE_SUBTYPE_CHECKOUT_PAYMENT:
				this.#paymentSelector.initPayments();
				show(this.#paymentSelector.htmlElement);
				break;
			default:
				throw `Unknown page type ${pageSubType}`;
				break;
		}
	}

	setOrder(order) {
		this.#checkoutAddress.setOrder(order);
		this.#shippingMethodSelector.setOrder(order);
		this.#paymentSelector.setOrder(order);
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
