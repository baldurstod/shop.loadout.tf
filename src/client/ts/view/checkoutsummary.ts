import { starSVG } from 'harmony-svg';
import { createElement, createShadowRoot, I18n, show } from 'harmony-ui';
import checkoutSummaryCSS from '../../css/checkoutsummary.css';
import commonCSS from '../../css/common.css';
import { Controller, ControllerEvent, NavigateToDetail } from '../controller';
import { Order } from '../model/order';
import { addressToHTML } from '../utils/address';
import { defineOrderSummary, HTMLOrderSummaryElement } from './components/ordersummary';
import { ShopElement } from './shopelement';
export { HTMLShopAddressElement } from './components/address';

export class CheckoutSummary extends ShopElement {
	#order?: Order;
	#htmlConfirm?: HTMLInputElement;
	#htmlContinue?: HTMLButtonElement;
	#htmlAddresses?: HTMLElement;
	#htmlSummary?: HTMLOrderSummaryElement;
	#htmlShippingAddress?: HTMLElement;
	#htmlBillingAddress?: HTMLElement;

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}

		defineOrderSummary();

		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [checkoutSummaryCSS, commonCSS],
			childs: [
				createElement('div', { class: 'title', i18n: '#order_summary' }),
				this.#htmlSummary = createElement('order-summary') as HTMLOrderSummaryElement,
				this.#htmlAddresses = createElement('div', {
					class: 'addresses',
					childs: [
						this.#htmlShippingAddress = createElement('div', {
							hidden: true,
						}) as HTMLElement,
						this.#htmlBillingAddress = createElement('div', {
							hidden: true,
						}) as HTMLElement,
					]
				}) as HTMLElement,
				createElement('div', {
					class: 'personalized-items',
					childs: [
						createElement('div', {
							innerHTML: starSVG,
						}),
						createElement('div', {
							class: 'personalized-items-text',
							childs: [
								createElement('h3', { i18n: '#personalized_items' }),
								createElement('span', { i18n: '#personalized_items_line1' }),
								createElement('br'),
								createElement('span', { i18n: '#personalized_items_line2' }),
							],
						}),
						createElement('img', {
							class: 'personalized-items-image',
							src: './img/custom.png',

						}),
					],
				}),
				createElement('button', {
					child: createElement('label', {
						class: 'confirm-label',
						childs: [
							createElement('span', {
								i18n: '#confirm_your_order',
							}),
							this.#htmlConfirm = createElement('input', {
								type: 'checkbox',
								$change: (event: InputEvent) => this.#confirmCheckout(event),
							}) as HTMLInputElement,
						],
					}),
				}),
				this.#htmlContinue = createElement('button', {
					i18n: '#continue_to_payment',
					disabled: true,
					$click: () => this.#continueCheckout(),
				}) as HTMLButtonElement,
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#refreshHTML(): void {
		if (!this.#order) {
			return;
		}

		this.initHTML();

		this.#htmlShippingAddress!.innerHTML = addressToHTML(this.#order.billingAddress, false);
		show(this.#htmlShippingAddress);
		this.#htmlBillingAddress!.innerHTML = addressToHTML(this.#order.billingAddress, true);
		show(this.#htmlBillingAddress);
	}

	setOrder(order: Order): void {
		this.initHTML();
		this.#order = order;
		this.#htmlSummary!.setOrder(order);
		this.#refreshHTML();

		// Reset the confirmation
		this.#htmlContinue!.disabled = true;
		this.#htmlConfirm!.checked = false;
	}

	#confirmCheckout(event: InputEvent): void {
		this.#htmlContinue!.disabled = !(event.target as HTMLInputElement).checked;
	}

	#continueCheckout(): void {
		Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@checkout#payment' } });
		/*
		if (this.#checkAddresses()) {
		}
		*/
	}
}
