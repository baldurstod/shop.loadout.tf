import { Address } from '../model/address';


export function addressToHTML(address: Address, billing: boolean): string {
	return `
		<div class="address">
			<div class="address-type i18n" data-i18n="${billing ? '#billing_address' : "#shipping_address"}"></div>
			<div class="address-line">${address.firstName} ${address.lastName}</div>
			<div class="address-line">${address.organization}</div>
			<div class="address-line">${address.address1}</div>
			<div class="address-line">${address.address2}</div>
			<div class="address-line">${address.city} ${address.postalCode}</div>
			<div class="address-line">${address.stateName} ${address.countryName}</div>
		</div>
	`;
}
