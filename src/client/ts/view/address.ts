import { createElement } from 'harmony-ui';
import { Address } from '../model/address';


export function address(address: Address, i18n: string): HTMLElement {
	const childs = [];

	if (i18n) {
		childs.push(createElement('div', { class: 'type', i18n: i18n, }));
	}

	childs.push(createElement('div', { class: 'name', innerText: `${address.firstName} ${address.lastName}`, }));
	childs.push(createElement('div', { class: 'address', innerText: address.address1, }));
	childs.push(createElement('div', { class: 'address', innerText: address.address2, }));
	if (address.stateCode) {
		childs.push(createElement('div', { class: 'city', innerText: `${address.city}, ${address.stateCode} ${address.postalCode}`, }));
	} else {
		childs.push(createElement('div', { class: 'city', innerText: `${address.postalCode} ${address.city}`, }));
	}
	childs.push(createElement('div', { class: 'city', innerText: address.countryName, }));
	childs.push(createElement('div', { class: 'email', innerText: address.email, }));
	childs.push(createElement('div', { class: 'phone', innerText: address.phone, }));

	return createElement('div', {
		class: 'address-block',
		childs: childs,
	});
}
