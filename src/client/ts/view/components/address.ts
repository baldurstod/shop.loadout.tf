import { I18n, createElement, display, hide, shadowRootStyle, show } from 'harmony-ui';
import addressCSS from '../../../css/address.css';
import commonCSS from '../../../css/common.css';
import { Address } from '../../model/address';
import { Countries } from '../../model/countries';

export class HTMLShopAddressElement extends HTMLElement {
	#shadowRoot!: ShadowRoot;
	#address = new Address();
	#htmlAddressType!: HTMLElement;
	#htmlFirstName!: HTMLInputElement;
	#htmlLastName!: HTMLInputElement;
	#htmlPhone!: HTMLInputElement;
	#htmlEmail!: HTMLInputElement;
	#htmlAddress1!: HTMLInputElement;
	#htmlAddress2!: HTMLInputElement;
	#htmlCountry!: HTMLSelectElement;
	#htmlState!: HTMLSelectElement;
	#htmlStateLine!: HTMLElement;
	#htmlPostalCode!: HTMLInputElement;
	#htmlCity!: HTMLInputElement;
	#countries?: Countries;
	#addressType = '';
	#htmlForm!: HTMLFormElement;

	constructor() {
		super();
		this.#initHTML();
	}

	#initHTML(): void {
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
		}) as HTMLFormElement;

		createElement('line', {
			parent: this.#htmlForm,
			childs: [
				createElement('label', {
					childs: [
						createElement('span', { i18n: '#first_name' }),
						this.#htmlFirstName = createElement('input', {
							i18n: { placeholder: '#first_name', },
							tabindex: 0,
							events: {
								input: (event: InputEvent) => this.#address.setFirstName((event.target as HTMLInputElement).value),
							}
						}) as HTMLInputElement,
					],
				}),

				createElement('label', {
					childs: [
						createElement('span', { i18n: '#last_name' }),
						this.#htmlLastName = createElement('input', {
							i18n: { placeholder: '#last_name', },
							tabindex: 0,
							events: {
								input: (event: InputEvent) => this.#address.setLastName((event.target as HTMLInputElement).value),
							}
						}) as HTMLInputElement,
					],
				}),
			]
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#phone' }),
				this.#htmlPhone = createElement('input', {
					i18n: { placeholder: '#phone', },
					tabindex: 0,
					events: {
						input: (event: InputEvent) => this.#address.setPhone((event.target as HTMLInputElement).value),
					}
				}) as HTMLInputElement,
			],
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#email' }),
				this.#htmlEmail = createElement('input', {
					i18n: { placeholder: '#email', },
					tabindex: 0,
					events: {
						input: (event: InputEvent) => this.#address.setEmail((event.target as HTMLInputElement).value),
					}
				}) as HTMLInputElement,
			],
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#address_line1' }),
				this.#htmlAddress1 = createElement('input', {
					i18n: { placeholder: '#address_line1', },
					tabindex: 0,
					events: {
						input: (event: InputEvent) => this.#address.setAddress1((event.target as HTMLInputElement).value),
					}
				}) as HTMLInputElement,
			],
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#address_line2' }),
				this.#htmlAddress2 = createElement('input', {
					i18n: { placeholder: '#address_line2', },
					tabindex: 0,
					events: {
						input: (event: InputEvent) => this.#address.setAddress2((event.target as HTMLInputElement).value),
					}
				}) as HTMLInputElement,
			],
		});

		createElement('label', {
			parent: this.#htmlForm,
			childs: [
				createElement('span', { i18n: '#country' }),
				this.#htmlCountry = createElement('select', {
					tabindex: 0,
					events: {
						input: (event: Event) => this.#selectCountry((event.target as HTMLSelectElement).value),
					}
				}) as HTMLSelectElement,
			],
		});

		this.#htmlStateLine = createElement('label', {
			parent: this.#htmlForm,
			hidden: true,
			childs: [
				createElement('span', { i18n: '#state' }),
				this.#htmlState = createElement('select', {
					tabindex: 0,
					events: {
						input: (event: Event) => this.#selectState((event.target as HTMLSelectElement).value),
					}
				}) as HTMLSelectElement,
			],
		});

		createElement('line', {
			parent: this.#htmlForm,
			childs: [
				createElement('label', {
					childs: [
						createElement('span', { i18n: '#postal_code' }),
						this.#htmlPostalCode = createElement('input', {
							i18n: { placeholder: '#postal_code', },
							tabindex: 0,
							events: {
								input: (event: InputEvent) => this.#address.setPostalCode((event.target as HTMLInputElement).value),
							}
						}) as HTMLInputElement,
					],
				}),

				createElement('label', {
					childs: [
						createElement('span', { i18n: '#city' }),
						this.#htmlCity = createElement('input', {
							i18n: { placeholder: '#city', },
							tabindex: 0,
							events: {
								input: (event: InputEvent) => this.#address.setCity((event.target as HTMLInputElement).value),
							}
						}) as HTMLInputElement,
					],
				}),
			]
		});
	}

	#refresh(): void {
		this.#htmlFirstName.value = this.#address.getFirstName();
		this.#htmlLastName.value = this.#address.getLastName();
		this.#htmlPhone.value = this.#address.getPhone();
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

				this.#htmlState.innerText = '';
				this.#htmlState.append(createElement('option'));
				for (const [, state] of country.getStates()) {
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

	setAddress(address: Address): void {
		this.#address = address;
		this.#refresh();
	}

	setCountries(countries: Countries): void {
		this.#countries = countries;
		console.info(countries);
		this.#htmlCountry.innerText = '';
		this.#htmlCountry.append(createElement('option'));

		for (const country of countries) {
			createElement('option', {
				parent: this.#htmlCountry,
				innerText: country.getName(),
				value: country.getCode(),
			});
		}

		this.#refresh();
	}

	#selectCountry(countryCode: string): void {
		this.#address.setCountryCode(countryCode);
		this.#address.setStateCode('');
		this.#refresh();
	}

	#selectState(stateCode: string): void {
		this.#address.setStateCode(stateCode);
		this.#refresh();
	}

	setAddressType(addressType: string): void {
		this.#addressType = addressType;
		this.#refresh();
	}

	checkAddress(): boolean {
		const address = this.#address;

		if (address.firstName == '') {
			this.#htmlFirstName.focus();
			return false;
		}

		if (address.lastName == '') {
			this.#htmlLastName.focus();
			return false;
		}

		if (address.address1 == '') {
			this.#htmlAddress1.focus();
			return false;
		}

		if (address.city == '') {
			this.#htmlCity.focus();
			return false;
		}

		if (address.postalCode == '') {
			this.#htmlPostalCode.focus();
			return false;
		}

		if (address.countryCode == '') {
			this.#htmlCountry.focus();
			return false;
		}

		if (address.phone == '') {
			this.#htmlPhone.focus();
			return false;
		}

		if (address.email == '') {
			this.#htmlEmail.focus();
			return false;
		}

		return true;

	}
}

let definedShopAddress = false;
export function defineShopAddress(): void {
	if (window.customElements && !definedShopAddress) {
		customElements.define('shop-address', HTMLShopAddressElement);
		definedShopAddress = true;
	}
}
