import { I18n, createElement, display } from 'harmony-ui';
import 'harmony-ui/dist/define/harmony-switch.js';
export { Address } from '../components/address.js';
import { PAYPAL_APP_CLIENT_ID } from '../../constants.js';
import { Controller } from '../../controller.js';
import { EVENT_NAVIGATE_TO } from '../../controllerevents.js';

import paypalCSS from '../../../css/payment/paypal.css';
import commonCSS from '../../../css/common.css';
import { fetchApi } from '../../fetchapi.js';


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

export class PaypalPayment {
	#htmlElement;
	#order;
	#paypalInitialized;

	constructor() {
		this.#initHTML();
	}

	async init(orderId) {
		if (this.#paypalInitialized) {
			return;
		}

		//await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&currency=${this.#order.currency}&intent=capture&enable-funding=venmo`)
		//await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&currency=USD&intent=capture&enable-funding=venmo`)
		await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&components=buttons&enable-funding=venmo`)
		console.info('paypal initialized')

		const paypalButtonsComponent = paypal.Buttons({
			// optional styling for buttons
			// https://developer.paypal.com/docs/checkout/standard/customize/buttons-style-guide/
				style: {
					color: "gold",
					shape: "rect",
					layout: "vertical"
				},

				// set up the transaction
				createOrder: async (data, actions) => {
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
					const approveResponse = await fetch('/paypal/order/capture', {
						method: 'POST',
						headers: {
							'Content-Type': 'application/json',
						},
						body: JSON.stringify({
							paypalOrderId: data.orderID,
						}),
					});
					const approveResponseJSON = await approveResponse.json();
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
			.render(this.paypalButtonContainer)
			.catch((err) => {
				console.error('PayPal Buttons failed to render');
			});

	}

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [ paypalCSS, commonCSS ],
			childs: [
				'paypal',


				this.paypalButtonContainer = createElement('div', {
					//id: 'paypal-button-container',
					id: 'test2',
				}),

			],
		});

		const test = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [ paypalCSS, commonCSS ],
			//style: 'display: none',
			parent: document.body,
			id: 'test',
			childs: [
				'paypal',
				this.paypalButtonContainer = createElement('div', {
					id: 'paypal-button-container',
				}),
			],
		}).host;

		setTimeout(() => document.body.append(test), 1000);
		//setTimeout(() => this.#htmlElement.append(test), 1000);

/*
		this.paypalButtonContainer = createElement('div', {
			parent: document.body
		}),
		*/

		I18n.observeElement(this.#htmlElement);
		this.init();
	}

	#refresh() {
	}

	setOrder(order) {
		this.#order = order;
		this.#refresh();
	}

	getHTMLElement() {
		return this.#htmlElement;
	}
}
