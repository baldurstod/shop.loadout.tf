import { HTMLHarmonySwitchElement, I18n, createElement, createShadowRoot, display } from 'harmony-ui';
export { HTMLShopAddressElement } from './components/address';
import { Controller } from '../controller';
import { EVENT_NAVIGATE_TO } from '../controllerevents';
import checkoutAddressesCSS from '../../css/checkoutaddresses.css';
import commonCSS from '../../css/common.css';
import { defineShopAddress, HTMLShopAddressElement } from './components/address';
import { Order } from '../model/order';
import { Countries } from '../model/countries';
import { ShopElement } from './shopelement';

export class CheckoutAddresses extends ShopElement {
	#htmlShippingAddress?: HTMLShopAddressElement;
	#htmlBillingAddress?: HTMLShopAddressElement;
	#htmlSameBillingAddress?: HTMLHarmonySwitchElement;
	#order?: Order;

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		defineShopAddress();
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [checkoutAddressesCSS, commonCSS],
			childs: [
				this.#htmlShippingAddress = createElement('shop-address', {
					elementCreated: (element: HTMLShopAddressElement) => element.setAddressType('#shipping_address'),
				}) as HTMLShopAddressElement,
				this.#htmlSameBillingAddress = createElement('harmony-switch', {
					'data-i18n': '#same_billing_address',
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
				}) as HTMLButtonElement,
			],
		});
		I18n.observeElement(this.shadowRoot);
	}

	#refreshHTML() {
		if (!this.#order) {
			return;
		}

		this.initHTML();

		const sameBillingAddress = this.#order?.getSameBillingAddress();
		this.#htmlSameBillingAddress!.state = sameBillingAddress;

		display(this.#htmlBillingAddress, !sameBillingAddress);
	}

	setOrder(order: Order) {
		this.initHTML();
		this.#order = order;
		this.#htmlShippingAddress!.setAddress(order.shippingAddress);
		this.#htmlBillingAddress!.setAddress(order.billingAddress);
		this.#refreshHTML();
	}

	setCountries(countries: Countries) {
		this.initHTML();
		this.#htmlShippingAddress!.setCountries(countries);
		this.#htmlBillingAddress!.setCountries(countries);
	}

	#changeSameBillingAddress(sameBillingAddress: boolean) {
		this.#order?.setSameBillingAddress(sameBillingAddress);
		this.#order && this.#htmlBillingAddress!.setAddress(this.#order.billingAddress);
		this.#refreshHTML();
	}

	#continueCheckout() {
		if (this.#checkAddresses()) {
			Controller.dispatchEvent(new CustomEvent(EVENT_NAVIGATE_TO, { detail: { url: '/@checkout#shipping' } }));
		}
	}

	#checkAddresses(): boolean {
		if  (!this.#htmlShippingAddress?.checkAddress()) {
			return false
		}

		if (!this.#htmlSameBillingAddress!.state) {
			if  (!this.#htmlBillingAddress?.checkAddress()) {
				return false
			}
		}

		return true;
	}
}
