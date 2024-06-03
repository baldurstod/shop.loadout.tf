import { I18n, createElement, display } from 'harmony-ui';
import 'harmony-ui/dist/define/harmony-switch.js';
export { Address } from '../components/address.js';

import paymentSelectorCSS from '../../../css/payment/paymentselector.css';
import commonCSS from '../../../css/common.css';

export class PaymentSelector {
	#htmlElement;
	#htmlMethods;
	#order;
	#payments = new Set();

	constructor() {
		this.#initHTML();
	}

	addPaymentMethod(payment) {
		this.#payments.add(payment);
		this.#refresh();
	}

	#initHTML() {
		console.info(this.#payments)

		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [ paymentSelectorCSS, commonCSS ],
			childs: [
				this.#htmlMethods = createElement('div', {
					class: 'payments',
					child: 'payment methods',
				}),
			],
		});
		I18n.observeElement(this.#htmlElement);
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
			this.#htmlMethods.append(payment.getHTMLElement());
		}
	}

	setOrder(order) {
		this.#order = order;
		this.#refresh();
	}

	get htmlElement() {
		return this.#htmlElement;
	}
}
