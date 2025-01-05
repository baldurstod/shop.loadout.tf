import { createElement } from 'harmony-ui';
import { Address } from '../model/address';


export function address(address: Address, i18n: string) {
	const childs = [];

	if (i18n) {
		childs.push(createElement('div', { class: 'type', i18n: i18n, }));
	}

	childs.push(createElement('div', { class: 'name', innerHTML: `${address.firstName} ${address.lastName}`, }));
	childs.push(createElement('div', { class: 'address', innerHTML: address.address1, }));
	childs.push(createElement('div', { class: 'address', innerHTML: address.address2, }));
	if (address.stateCode) {
		childs.push(createElement('div', { class: 'city', innerHTML: `${address.city}, ${address.stateCode} ${address.postalCode}`, }));
	} else {
		childs.push(createElement('div', { class: 'city', innerHTML: `${address.postalCode} ${address.city}`, }));
	}
	childs.push(createElement('div', { class: 'city', innerHTML: address.countryName, }));
	childs.push(createElement('div', { class: 'email', innerHTML: address.email, }));
	childs.push(createElement('div', { class: 'phone', innerHTML: address.phone, }));

	return createElement('div', {
		class: 'address-block',
		childs: childs,
	});
}
