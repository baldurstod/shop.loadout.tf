import { createElement, createShadowRoot, display, HTMLHarmonySwitchElement, I18n } from 'harmony-ui';
import checkoutAddressesCSS from '../../css/checkoutaddresses.css';
import commonCSS from '../../css/common.css';
import { Controller, ControllerEvent, NavigateToDetail } from '../controller';
import { Countries } from '../model/countries';
import { Order } from '../model/order';
import { defineShopAddress, HTMLShopAddressElement } from './components/address';
import { ShopElement } from './shopelement';
export { HTMLShopAddressElement } from './components/address';

export class CheckoutAddresses extends ShopElement {
	#htmlShippingAddress?: HTMLShopAddressElement;
	#htmlBillingAddress?: HTMLShopAddressElement;
	#htmlSameBillingAddress?: HTMLHarmonySwitchElement;
	#order?: Order;

	initHTML(): void {
		if (this.shadowRoot) {
			return;
		}
		defineShopAddress();
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [checkoutAddressesCSS, commonCSS],
			childs: [
				this.#htmlShippingAddress = createElement('shop-address', {
					elementCreated: (element: Element) => (element as HTMLShopAddressElement).setAddressType('#shipping_address'),
				}) as HTMLShopAddressElement,
				this.#htmlSameBillingAddress = createElement('harmony-switch', {
					'data-i18n': '#same_billing_address',
					events: {
						change: (event: CustomEvent) => this.#changeSameBillingAddress((event.target as HTMLHarmonySwitchElement).state ?? true),
					},
				}) as HTMLHarmonySwitchElement,
				this.#htmlBillingAddress = createElement('shop-address', {
					elementCreated: (element: Element) => (element as HTMLShopAddressElement).setAddressType('#billing_address'),
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

	#refreshHTML(): void {
		if (!this.#order) {
			return;
		}

		this.initHTML();

		const sameBillingAddress = this.#order?.getSameBillingAddress();
		this.#htmlSameBillingAddress!.state = sameBillingAddress;

		display(this.#htmlBillingAddress, !sameBillingAddress);
	}

	setOrder(order: Order): void {
		this.initHTML();
		this.#order = order;
		this.#htmlShippingAddress!.setAddress(order.shippingAddress);
		this.#htmlBillingAddress!.setAddress(order.billingAddress);
		this.#refreshHTML();
	}

	setCountries(countries: Countries): void {
		this.initHTML();
		this.#htmlShippingAddress!.setCountries(countries);
		this.#htmlBillingAddress!.setCountries(countries);
	}

	#changeSameBillingAddress(sameBillingAddress: boolean): void {
		this.#order?.setSameBillingAddress(sameBillingAddress);
		if (this.#order) {
			this.#htmlBillingAddress!.setAddress(this.#order.billingAddress);
		}
		this.#refreshHTML();
	}

	#continueCheckout(): void {
		if (this.#checkAddresses()) {
			Controller.dispatchEvent<NavigateToDetail>(ControllerEvent.NavigateTo, { detail: { url: '/@checkout#shipping' } });
		}
	}

	#checkAddresses(): boolean {
		if (!this.#htmlShippingAddress?.checkAddress()) {
			return false
		}

		if (!this.#htmlSameBillingAddress!.state) {
			if (!this.#htmlBillingAddress?.checkAddress()) {
				return false
			}
		}

		return true;
	}
}
