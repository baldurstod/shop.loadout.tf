import { I18n, createElement, createShadowRoot, display } from 'harmony-ui';
import { Controller } from '../controller';
import { EVENT_NAVIGATE_TO } from '../controllerevents';

import shippingMethodSelectorCSS from '../../css/shippingmethodselector.css';
import commonCSS from '../../css/common.css';
import { Order } from '../model/order';
import { ShopElement } from './shopelement';

export class ShippingMethodSelector extends ShopElement {
	#htmlMethods?: HTMLElement;
	#htmlContinue?: HTMLButtonElement;

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [shippingMethodSelectorCSS, commonCSS],
			childs: [
				this.#htmlMethods = createElement('div', {
					class: 'methods',
				}),
				this.#htmlContinue = createElement('button', {
					i18n: '#continue_to_payment',
					disabled: true,
					events: {
						click: () => this.#continueCheckout(),
					},
				}) as HTMLButtonElement,
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#refreshHTML(order: Order) {
		this.initHTML();

		this.#htmlMethods!.replaceChildren();
		console.info(order);
		console.info(order.shippingInfos);

		let htmlRadio: HTMLInputElement;
		for (const [_, shippingInfo] of order.shippingInfos) {
			createElement('label', {
				parent: this.#htmlMethods,
				class: 'method',
				childs: [
					htmlRadio = createElement('input', {
						type: 'radio',
						name: 'shipping-method',
						events: {
							input: (event: InputEvent) => {
								if ((event.target as HTMLInputElement).checked) {
									order.shippingMethod = shippingInfo.shipping;
								}
							}
						},
					}) as HTMLInputElement,
					createElement('div', {
						class: 'method-name',
						innerText: shippingInfo.shippingMethodName,
					}),
					createElement('div', {
						class: 'method-rate',
						innerText: shippingInfo.rate,
					}),
					createElement('div', {
						class: 'tick',
					}),
				]
			});

			if (shippingInfo.shipping == order.shippingMethod) {
				htmlRadio.checked = true;
			}
		}

		if (order.shippingInfos.size) {
			this.#htmlContinue?.removeAttribute('disabled');
		} else {
			this.#htmlContinue?.setAttribute('disabled', '1');
		}
	}

	setOrder(order: Order) {
		this.#refreshHTML(order);
	}

	#continueCheckout() {
		//TODO: check values
		Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout#payment' } }));
	}
}
