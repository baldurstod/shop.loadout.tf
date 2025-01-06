import { I18n, createElement, createShadowRoot, display } from 'harmony-ui';
import { Payment } from './payment';
import paymentSelectorCSS from '../../../css/payment/paymentselector.css';
import commonCSS from '../../../css/common.css';
import { Order } from '../../model/order';
import { ShopElement } from '../shopelement';

export class PaymentSelector extends ShopElement {
	#htmlMethods?: HTMLElement;
	#order?: Order;
	#payments = new Set<Payment>();

	addPaymentMethod(payment: Payment) {
		this.#payments.add(payment);
		this.#refreshHTML();
	}

	async initPayments() {
		for (const payment of this.#payments) {
			await payment.initPayment(this.#order?.id ?? '');
		}
	}

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		console.info(this.#payments)

		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [paymentSelectorCSS, commonCSS],
			childs: [
				'payment methods',
				this.#htmlMethods = createElement('div', {
					class: 'payments',
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#refreshHTML() {
		if (!this.#order) {
			return;
		}
		this.initHTML();

		this.#htmlMethods!.replaceChildren();
		console.info(this.#order);
		console.info(this.#order.shippingInfos);

		let htmlRadio;

		for (const payment of this.#payments) {
			this.#htmlMethods!.append(payment.getHTML());
		}
	}

	setOrder(order: Order) {
		this.#order = order;
		this.#refreshHTML();
	}
}
