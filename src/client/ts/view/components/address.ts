import { I18n, createElement, display, hide, shadowRootStyle, show } from 'harmony-ui';

import addressCSS from '../../../css/address.css';
import commonCSS from '../../../css/common.css';
import { Address } from '../../model/address';

export class HTMLShopAddressElement extends HTMLElement {
	#shadowRoot;
	#address = new Address();
	#htmlAddressType;
	#htmlFirstName;
	#htmlLastName;
	#htmlEmail;
	#htmlAddress1;
	#htmlAddress2;
	#htmlCountry;
	#htmlState;
	#htmlStateLine;
	#htmlPostalCode;
	#htmlCity;
	#countries;
	#addressType = '';
	#htmlForm;
	constructor() {
		super();
		this.#initHTML();
	}

	#initHTML() {
		this.#shadowRoot = this.attachShadow({ mode: 'closed' });
		I18n.observeElement(this.#shadowRoot);
		shadowRootStyle(this.#shadowRoot, addressCSS);
		shadowRootStyle(this.#shadowRoot, commonCSS);

		this.#htmlAddressType = createElement('h2', {
			parent: this.#shadowRoot,
			i18n: '',
		});

		this.#htmlForm = createElement('form', {
			parent: this.#shadowRoot
		});

		createElement('line', {
			parent: this.#htmlForm,
			childs: [
				createElement('label', {
					childs: [
						createElement('span', { i18n: '#first_name' }),
						this.#htmlFirstName = createElement('input', {
							'i18n-placeholder': '#first_name',
							events: {
								input: event => this.#address.setFirstName(event.target.value),
							}
						}),
					],
				}),

				createElement('label', {
					childs: [
						createElement('span', { i18n: '#last_name' }),
						this.#htmlLastName = createElement('input', {
							'i18n-placeholder': '#last_name',
							events: {
								input: event => this.#address.setLastName(event.target.value),
							}
						}),
					],
				}),
			]
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#email' }),
				this.#htmlEmail = createElement('input', {
					'i18n-placeholder': '#email',
					events: {
						input: event => this.#address.setEmail(event.target.value),
					}
				}),
			],
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#address_line1' }),
				this.#htmlAddress1 = createElement('input', {
					'i18n-placeholder': '#address_line1',
					events: {
						input: event => this.#address.setAddress1(event.target.value),
					}
				}),
			],
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#address_line2' }),
				this.#htmlAddress2 = createElement('input', {
					'i18n-placeholder': '#address_line2',
					events: {
						input: event => this.#address.setAddress2(event.target.value),
					}
				}),
			],
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#country' }),
				this.#htmlCountry = createElement('select', {
					events: {
						input: event => this.#selectCountry(event.target.value),
					}
				}),
			],
		});

		this.#htmlStateLine = createElement('label', {
			parent: this.#htmlForm,
			hidden: true,
			childs: [
				createElement('span', { i18n: '#state' }),
				this.#htmlState = createElement('select', {
					events: {
						input: event => this.#selectState(event.target.value),
					}
				}),
			],
		});

		createElement('line', {
			parent: this.#htmlForm,
			childs: [
				createElement('label', {
					childs: [
						createElement('span', { i18n: '#postal_code' }),
						this.#htmlPostalCode = createElement('input', {
							'i18n-placeholder': '#postal_code',
							events: {
								input: event => this.#address.setPostalCode(event.target.value),
							}
						}),
					],
				}),

				createElement('label', {
					childs: [
						createElement('span', { i18n: '#city' }),
						this.#htmlCity = createElement('input', {
							'i18n-placeholder': '#city',
							events: {
								input: event => this.#address.setCity(event.target.value),
							}
						}),
					],
				}),
			]
		});
	}

	#refresh() {
		this.#htmlFirstName.value = this.#address.getFirstName();
		this.#htmlLastName.value = this.#address.getLastName();
		this.#htmlEmail.value = this.#address.getEmail();
		this.#htmlAddress1.value = this.#address.getAddress1();
		this.#htmlAddress2.value = this.#address.getAddress2();
		this.#htmlPostalCode.value = this.#address.getPostalCode();
		this.#htmlCity.value = this.#address.getCity();

		const countryCode = this.#address.getCountryCode();
		const country = this.#countries?.getCountry(countryCode);
		if (country) {
			this.#htmlCountry.value = countryCode;
			console.info(country);
			if (country.hasStates()) {
				show(this.#htmlStateLine);

				this.#htmlState.innerHTML = '';
				this.#htmlState.append(createElement('option'));
				for (let [_, state] of country.getStates()) {
					createElement('option', {
						parent: this.#htmlState,
						innerText: state.getName(),
						value: state.getCode(),
					});
				}

				this.#htmlState.value = this.#address.getStateCode();
			} else {
				hide(this.#htmlStateLine);
			}
		} else {
			this.#htmlCountry.value = '';
			this.#htmlState.value = '';
		}

		display(this.#htmlAddressType, this.#addressType != '');
		this.#htmlAddressType.setAttribute('data-i18n', this.#addressType);
	}

	setAddress(address) {
		this.#address = address;
		this.#refresh();
	}

	setCountries(countries = []) {
		this.#countries = countries;
		console.info(countries);
		this.#htmlCountry.innerHTML = '';
		this.#htmlCountry.append(createElement('option'));

		for (let country of countries) {
			createElement('option', {
				parent: this.#htmlCountry,
				innerText: country.getName(),
				value: country.getCode(),
			});
		}

		this.#refresh();
	}

	#selectCountry(countryCode) {
		this.#address.setCountryCode(countryCode);
		this.#address.setStateCode('');
		this.#refresh();
	}

	#selectState(stateCode) {
		this.#address.setStateCode(stateCode);
		this.#refresh();
	}

	setAddressType(addressType) {
		this.#addressType = addressType;
		this.#refresh();
	}
}

let definedShopAddress = false;
export function defineShopAddress() {
	if (window.customElements && !definedShopAddress) {
		customElements.define('shop-address', HTMLShopAddressElement);
		definedShopAddress = true;
	}
}
