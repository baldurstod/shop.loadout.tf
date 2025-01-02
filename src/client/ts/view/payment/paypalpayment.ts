import { I18n, createElement, display } from 'harmony-ui';
import { PAYPAL_APP_CLIENT_ID } from '../../constants';
import { Controller } from '../../controller';
import { EVENT_NAVIGATE_TO } from '../../controllerevents';
import { Payment } from './payment';
import { fetchApi } from '../../fetchapi';

import paypalCSS from '../../../css/payment/paypal.css';
import paypalButtonsCSS from '../../../css/payment/paypalbuttons.css';
import commonCSS from '../../../css/common.css';


export function loadScript(scriptSrc) {
	return new Promise((resolve) => {
		const script = createElement('script', {
			src: scriptSrc,
			parent: document.head,
			events: {
				load: () => resolve(true),
			}
		});
	});
}

export class PaypalPayment extends Payment {
	#htmlElement;
	#order;
	#paypalInitialized;
	#paypalDialog;
	#paypalButtonContainer;

	constructor() {
		super();
		this.#initHTML();
	}

	async initPayment(orderId) {
		if (this.#paypalInitialized) {
			return;
		}

		//await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&currency=${this.#order.currency}&intent=capture&enable-funding=venmo`)
		//await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&currency=USD&intent=capture&enable-funding=venmo`)
		await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&components=buttons&enable-funding=venmo`)
		console.info('paypal initialized')

		const paypalButtonsComponent = (window as any).paypal.Buttons({
			// optional styling for buttons
			// https://developer.paypal.com/docs/checkout/standard/customize/buttons-style-guide/
			style: {
				color: "gold",
				shape: "rect",
				layout: "vertical"
			},

			// set up the transaction
			createOrder: async (data, actions) => {
				this.#paypalDialog.close();
				const { requestId, response } = await fetchApi({
					action: 'create-paypal-order',
					version: 1,
					params: {
						order_id: orderId,
					},
				});

				if (response.success) {
					return response.result.paypal_order_id;
				} else {
					console.error('Error while creating paypal order', response);
					throw 'Something wrong happened';
				}


				/*const response = await fetch('/paypal/order/create', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						id: orderId,
					}),
				});*/
			},

			// finalize the transaction
			onApprove: async (data, actions) => {
				/*
				const approveResponse = await fetch('/paypal/order/capture', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						paypalOrderId: data.orderID,
					}),
				});
				*/
				const { requestId, response } = await fetchApi({
					action: 'capture-paypal-order',
					version: 1,
					params: {
						paypal_order_id: data.orderID,
					},
				});

				const approveResponseJSON = await response.json();
				if (approveResponseJSON.success) {
					Controller.dispatchEvent(new CustomEvent('paymentcomplete', { detail: approveResponseJSON.order }));
				}
			},

			// handle unrecoverable errors
			onError: (err) => {
				console.error('An error prevented the buyer from checking out with PayPal');
			}
		});

		paypalButtonsComponent
			.render(this.#paypalButtonContainer)
			.catch((err) => {
				console.error('PayPal Buttons failed to render');
			});

		this.#paypalDialog.showModal()
	}

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [paypalCSS, commonCSS],
			childs: [
				createElement('div', {
					i18n: '#select_paypal_payment',
				}),
			],
			events: {
				click: () => this.#paypalDialog.showModal(),

			},
		});

		this.#paypalDialog = createElement('dialog', {
			//attachShadow: { mode: 'closed' },
			//adoptStyles: [ paypalButtonsCSS ],
			parent: document.body,
			class: 'paypal-dialog',
			childs: [
				this.#paypalButtonContainer = createElement('div', {
					id: 'paypal-button-container',
				}),
			],
		});

		I18n.observeElement(this.#htmlElement);
	}

	#refresh() {
	}

	setOrder(order) {
		this.#order = order;
		this.#refresh();
	}

	getHtmlElement() {
		return this.#htmlElement.host;
	}
}
