import { I18n, createElement, createShadowRoot, display } from 'harmony-ui';
import { Controller } from '../controller';
import { EVENT_NAVIGATE_TO } from '../controllerevents';

import shippingMethodSelectorCSS from '../../css/shippingmethodselector.css';
import commonCSS from '../../css/common.css';

export class ShippingMethodSelector {
	#shadowRoot?: ShadowRoot;
	#htmlMethods;
	#order;

	constructor() {
		this.#initHTML();
	}

	#initHTML() {
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyles: [shippingMethodSelectorCSS, commonCSS],
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

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}
}
