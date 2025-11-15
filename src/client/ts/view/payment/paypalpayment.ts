import { addNotification, NotificationType } from 'harmony-browser-utils';
import { createElement, createShadowRoot, I18n } from 'harmony-ui';
import commonCSS from '../../../css/common.css';
import paypalCSS from '../../../css/payment/paypal.css';
import { PAYPAL_APP_CLIENT_ID } from '../../constants';
import { Controller, ControllerEvent, PaymentCancelledDetail } from '../../controller';
import { fetchApi } from '../../fetchapi';
import { Order } from '../../model/order';
import { CreatePaypalOrderResponse } from '../../responses/createpaypalorder';
import { CapturePaypalOrderResponse, OrderJSON } from '../../responses/order';
import { ShopElement } from '../shopelement';
import { Payment } from './payment';


export function loadScript(scriptSrc: string): Promise<boolean> {
	return new Promise<boolean>((resolve) => {
		createElement('script', {
			src: scriptSrc,
			parent: document.head,
			events: {
				load: () => resolve(true),
			}
		});
	});
}

type PaypalData = {
	orderID: string,
}

type PaypalButton = {
	render: (arg: HTMLElement) => PaypalButton;
	catch: (err: unknown) => void;
}

type PaypalWindow = Window & typeof globalThis & {
	paypal: {
		Buttons: (arg1: object) => PaypalButton;
	}
}

export class PaypalPayment extends ShopElement implements Payment {
	isPayment = true as const;
	#paypalInitialized = false;
	#paypalDialog?: HTMLDialogElement;
	#paypalButtonContainer?: HTMLElement;

	async initPayment(): Promise<void> {
		if (this.#paypalInitialized) {
			return;
		}

		this.initHTML();

		//await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&currency=${this.#order.currency}&intent=capture&enable-funding=venmo`)
		//await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&currency=USD&intent=capture&enable-funding=venmo`)
		await loadScript(`https://www.paypal.com/sdk/js?client-id=${PAYPAL_APP_CLIENT_ID}&components=buttons&enable-funding=venmo`)
		console.info('paypal initialized')

		const paypalButtonsComponent: PaypalButton = (window as PaypalWindow).paypal.Buttons({
			// optional styling for buttons
			// https://developer.paypal.com/docs/checkout/standard/customize/buttons-style-guide/
			style: {
				color: 'gold',
				shape: 'rect',
				layout: 'vertical'
			},

			// set up the transaction
			createOrder: async (/*data, actions*/) => {
				this.#paypalDialog!.close();
				const { requestId, response } = await fetchApi('create-paypal-order', 1) as { requestId: string, response: CreatePaypalOrderResponse };

				if (response.success) {
					return response.result?.paypal_order_id;
				} else {
					console.error('Error while creating paypal order', response, requestId);
					throw new Error('Something wrong happened');
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
			onApprove: async (data: PaypalData, /*actions*/) => {
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
				const { requestId, response } = await fetchApi('capture-paypal-order', 1, {
					paypal_order_id: data.orderID,
				}) as { requestId: string, response: CapturePaypalOrderResponse };

				if (response.success && response.result) {
					Controller.dispatchEvent<OrderJSON>(ControllerEvent.PaymentComplete, { detail: response.result.order });
				} else {
					addNotification(createElement('span', {
						i18n: {
							innerText: '#error_while_processing_payment',
							values: {
								requestId: requestId,
							},
						},
					}), NotificationType.Error, 0);
				}

			},

			onCancel: function (data: PaypalData) {
				Controller.dispatchEvent<PaymentCancelledDetail>(ControllerEvent.PaymentCancelled, { detail: { orderId: data.orderID } });
			},

			// handle unrecoverable errors
			onError: (err: unknown) => {
				console.error('An error prevented the buyer from checking out with PayPal', err);
			}
		});

		paypalButtonsComponent
			.render(this.#paypalButtonContainer!)
			.catch((err: unknown) => {
				console.error('PayPal Buttons failed to render', err);
			});

		this.#paypalDialog!.showModal()
	}

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [paypalCSS, commonCSS],
			childs: [
				createElement('div', {
					i18n: '#select_paypal_payment',
				}),
			],
			events: {
				click: () => this.#paypalDialog!.showModal(),

			},
		});

		this.#paypalDialog = createElement('dialog', {
			parent: document.body,
			class: 'paypal-dialog',
			childs: [
				this.#paypalButtonContainer = createElement('div', {
					id: 'paypal-button-container',
				}),
			],
		}) as HTMLDialogElement;
		I18n.observeElement(this.shadowRoot);
	}

	#refreshHTML(order: Order): void {
		console.error(order);
	}

	setOrder(order: Order): void {
		this.#refreshHTML(order);
	}
}
