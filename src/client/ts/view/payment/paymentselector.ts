import { I18n, createElement, createShadowRoot, display } from 'harmony-ui';
import { Payment } from './payment';

import paymentSelectorCSS from '../../../css/payment/paymentselector.css';
import commonCSS from '../../../css/common.css';

export class PaymentSelector {
	#shadowRoot?: ShadowRoot;
	#htmlMethods;
	#order;
	#payments = new Set<Payment>();

	constructor() {
		this.#initHTML();
	}

	addPaymentMethod(payment) {
		this.#payments.add(payment);
		this.#refresh();
	}

	async initPayments() {
		for (const payment of this.#payments) {
			await payment.initPayment(null);
		}
	}

	#initHTML() {
		console.info(this.#payments)

		this.#shadowRoot = createShadowRoot('section', {
			adoptStyles: [ paymentSelectorCSS, commonCSS ],
			childs: [
				'payment methods',
				this.#htmlMethods = createElement('div', {
					class: 'payments',
				}),
			],
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	#refresh() {
		if (!this.#order) {
			return;
		}

		this.#htmlMethods.replaceChildren();
		console.info(this.#order);
		console.info(this.#order.shippingInfos);

		let htmlRadio;

		for (const payment of this.#payments) {
			this.#htmlMethods.append(payment.getHtmlElement());
		}
	}

	setOrder(order) {
		this.#order = order;
		this.#refresh();
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}
}
