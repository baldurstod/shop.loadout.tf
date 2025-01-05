import { Address } from './model/address';

/*export function addressFromPaypalUserInfo(userInfo) {
	const address = userInfo.address;

	return {
		name: userInfo.name,
		email: userInfo.email,

		address1: address?.street_address,
		city: address?.locality,
		stateCode: address?.region,
		countryCode: address?.country,
		postalCode: address?.postal_code,

		verified: (userInfo?.verified === true) || (userInfo?.verified === 'true'),
	};
}*/

export function addressToPaypalShipping(address: Address) {
	return {
		name: {
			full_name: address.name
		},
		address: {
			address_line_1: address.address1,
			admin_area_1: address.stateCode,
			admin_area_2: address.city,
			postal_code: address.postalCode,
			country_code: address.countryCode

		},
		type: 'SHIPPING',
	};
}

export function addressToPrintfulRecipient(address: Address) {
	return {
		name: address.name,
		email: address.email,

		address1: address.address1,
		city: address.city,
		zip: address.postalCode,
		country_code: address.countryCode,
		state_code: address.stateCode,
	};
}
