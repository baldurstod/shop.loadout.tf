import { createShadowRoot, hide, I18n, show } from 'harmony-ui';
import { CheckoutAddresses } from './checkoutaddresses';
import { PaymentSelector } from './payment/paymentselector';
import { ShippingMethodSelector } from './shippingmethodselector';
import { PaypalPayment } from './payment/paypalpayment';

import checkoutPageCSS from '../../css/checkoutpage.css';
import { PAGE_SUBTYPE_CHECKOUT_ADDRESS, PAGE_SUBTYPE_CHECKOUT_INIT, PAGE_SUBTYPE_CHECKOUT_PAYMENT, PAGE_SUBTYPE_CHECKOUT_SHIPPING, PageSubType } from '../constants.js';
import { Order } from '../model/order';
import { Countries } from '../model/countries';

export class CheckoutPage {
	#shadowRoot?: ShadowRoot;
	#checkoutAddress = new CheckoutAddresses();
	#shippingMethodSelector = new ShippingMethodSelector();
	#paymentSelector = new PaymentSelector();

	constructor() {
		this.#initHTML();

		this.#paymentSelector.addPaymentMethod(new PaypalPayment());
	}

	#initHTML() {
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyle: checkoutPageCSS,
			childs: [
				this.#checkoutAddress.getHTML(),
				this.#shippingMethodSelector.getHTML(),
				this.#paymentSelector.getHTML(),
			],
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	setCheckoutStage(pageSubType: PageSubType) {
		hide(this.#checkoutAddress.getHTML());
		hide(this.#shippingMethodSelector.getHTML());
		hide(this.#paymentSelector.getHTML());
		switch (pageSubType) {
			case PAGE_SUBTYPE_CHECKOUT_INIT:
				break;
			case PAGE_SUBTYPE_CHECKOUT_ADDRESS:
				show(this.#checkoutAddress.getHTML());
				break;
			case PAGE_SUBTYPE_CHECKOUT_SHIPPING:
				show(this.#shippingMethodSelector.getHTML());
				break;
			case PAGE_SUBTYPE_CHECKOUT_PAYMENT:
				this.#paymentSelector.initPayments();
				show(this.#paymentSelector.getHTML());
				break;
			default:
				throw `Unknown page type ${pageSubType}`;
				break;
		}
	}

	setOrder(order: Order) {
		this.#checkoutAddress.setOrder(order);
		this.#shippingMethodSelector.setOrder(order);
		this.#paymentSelector.setOrder(order);
	}

	setCountries(countries: Countries) {
		this.#checkoutAddress.setCountries(countries);
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}
}
