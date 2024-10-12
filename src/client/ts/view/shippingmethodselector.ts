import { I18n, createElement, display } from 'harmony-ui';
import 'harmony-ui/dist/define/harmony-switch.js';
export { Address } from './components/address.js';
import { Controller } from '../controller.js';
import { EVENT_NAVIGATE_TO } from '../controllerevents.js';

import shippingMethodSelectorCSS from '../../css/shippingmethodselector.css';
import commonCSS from '../../css/common.css';

export class ShippingMethodSelector {
	#htmlElement;
	#htmlMethods;
	#order;

	constructor() {
		this.#initHTML();
	}

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [ shippingMethodSelectorCSS, commonCSS ],
			childs: [
				this.#htmlMethods = createElement('div', {
					class: 'methods',
				}),
				createElement('button', {
					i18n: '#continue_to_payment',
					events: {
						click: () => this.#continueCheckout(),
					},
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
		for (const [_, shippingInfo] of this.#order.shippingInfos) {
			createElement('label', {
				parent: this.#htmlMethods,
				class: 'method',
				childs: [
					htmlRadio = createElement('input', {
						type: 'radio',
						name: 'shipping-method',
						events: {
							input: (event) => {
								if (event.target.checked) {
									this.#order.shippingMethod = shippingInfo.id;
								}
							}
						},
					}),
					createElement('div', {
						class: 'method-name',
						innerText: shippingInfo.name,
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

			if (shippingInfo.id == this.#order.shippingMethod) {
				htmlRadio.checked = true;
			}
		}
	}

	setOrder(order) {
		this.#order = order;
		this.#refresh();
	}

	#continueCheckout() {
		//TODO: check values
		Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout#payment' } }));
	}

	get htmlElement() {
		return this.#htmlElement;
	}
}
