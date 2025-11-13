import { I18n, createElement, createShadowRoot } from 'harmony-ui';
import commonCSS from '../../../css/common.css';
import paymentSelectorCSS from '../../../css/payment/paymentselector.css';
import { Order } from '../../model/order';
import { ShopElement } from '../shopelement';
import { Payment } from './payment';

export class PaymentSelector extends ShopElement {
	#htmlMethods?: HTMLElement;
	#order?: Order;
	#payments = new Set<Payment>();

	addPaymentMethod(payment: Payment): void {
		this.#payments.add(payment);
		this.#refreshHTML();
	}

	async initPayments(): Promise<void> {
		for (const payment of this.#payments) {
			await payment.initPayment();
		}
	}

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}
		console.info(this.#payments)

		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [paymentSelectorCSS, commonCSS],
			childs: [
				this.#htmlMethods = createElement('div', {
					i18n: '#payment_methods',
				}),
				this.#htmlMethods = createElement('div', {
					class: 'payments',
				}),
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#refreshHTML(): void {
		if (!this.#order) {
			return;
		}
		this.initHTML();

		this.#htmlMethods!.replaceChildren();
		console.info(this.#order);
		console.info(this.#order.shippingInfos);

		for (const payment of this.#payments) {
			this.#htmlMethods!.append(payment.getHTML());
		}
	}

	setOrder(order: Order): void {
		this.#order = order;
		this.#refreshHTML();
	}
}
