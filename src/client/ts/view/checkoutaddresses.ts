import { HTMLHarmonySwitchElement, I18n, createElement, createShadowRoot, display } from 'harmony-ui';
export { HTMLShopAddressElement } from './components/address';
import { Controller } from '../controller';
import { EVENT_NAVIGATE_TO } from '../controllerevents';
import checkoutAddressesCSS from '../../css/checkoutaddresses.css';
import commonCSS from '../../css/common.css';
import { defineShopAddress, HTMLShopAddressElement } from './components/address';
import { Order } from '../model/order';

export class CheckoutAddresses {
	#shadowRoot?: ShadowRoot;
	#htmlShippingAddress?: HTMLShopAddressElement;
	#htmlBillingAddress?: HTMLShopAddressElement;
	#htmlSameBillingAddress?: HTMLHarmonySwitchElement;
	#order?: Order;

	#initHTML() {
		defineShopAddress();
		this.#shadowRoot = createShadowRoot('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [checkoutAddressesCSS, commonCSS],
			childs: [
				this.#htmlShippingAddress = createElement('shop-address', {
					elementCreated: (element: HTMLShopAddressElement) => element.setAddressType('#shipping_address'),
				}) as HTMLShopAddressElement,
				this.#htmlSameBillingAddress = createElement('harmony-switch', {
					i18n: '#same_billing_address',
					events: {
						change: (event: CustomEvent) => this.#changeSameBillingAddress((event.target as HTMLHarmonySwitchElement).state ?? true),
					},
				}) as HTMLHarmonySwitchElement,
				this.#htmlBillingAddress = createElement('shop-address', {
					elementCreated: (element: HTMLShopAddressElement) => element.setAddressType('#billing_address'),
				}) as HTMLShopAddressElement,
				createElement('button', {
					i18n: '#continue_to_shipping',
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

		if (!this.#shadowRoot) {
			this.#initHTML();
		}

		const sameBillingAddress = this.#order?.getSameBillingAddress();
		this.#htmlSameBillingAddress!.state = sameBillingAddress;

		display(this.#htmlBillingAddress, !sameBillingAddress);
	}

	setOrder(order: Order) {
		if (!this.#shadowRoot) {
			this.#initHTML();
		}
		this.#order = order;
		this.#htmlShippingAddress!.setAddress(order.shippingAddress);
		this.#htmlBillingAddress!.setAddress(order.billingAddress);
		this.#refresh();
	}

	setCountries(countries) {
		if (!this.#shadowRoot) {
			this.#initHTML();
		}
		this.#htmlShippingAddress!.setCountries(countries);
		this.#htmlBillingAddress!.setCountries(countries);
	}

	#changeSameBillingAddress(sameBillingAddress: boolean) {
		this.#order?.setSameBillingAddress(sameBillingAddress);
		this.#refresh();
	}

	#continueCheckout() {
		//TODO: check values
		Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout#shipping' } }));
	}

	get htmlElement() {
		throw 'use getHTML';
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}
}
