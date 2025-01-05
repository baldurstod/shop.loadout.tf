import { I18n, createElement, createShadowRoot, display } from 'harmony-ui';
import { Payment } from './payment';
import paymentSelectorCSS from '../../../css/payment/paymentselector.css';
import commonCSS from '../../../css/common.css';
import { Order } from '../../model/order';

export class PaymentSelector {
	#shadowRoot?: ShadowRoot;
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

	#initHTML() {
		console.info(this.#payments)

		this.#shadowRoot = createShadowRoot('section', {
			adoptStyles: [paymentSelectorCSS, commonCSS],
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

	#refreshHTML() {
		if (!this.#order) {
			return;
		}
		if (!this.#shadowRoot) {
			this.#initHTML();
		}

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

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}
}
