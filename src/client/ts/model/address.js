export class Address {
	constructor() {
		this.firstName = '';
		this.lastName = '';
		this.organization = '';
		this.address1 = '';
		this.address2 = '';
		this.city = '';
		this.stateCode = '';
		this.stateName = '';
		this.countryCode = '';
		this.countryName = '';
		this.postalCode = '';
		this.phone = '';
		this.email = '';
		this.taxNumber = '';
	}

	get name() {
		return `${this.firstName} ${this.lastName}`;
	}

	fromJSON(json) {
		if (!json) {
			return;
		}
		this.firstName = json.first_name;
		this.lastName = json.last_name;
		this.organization = json.organization;
		this.address1 = json.address1;
		this.address2 = json.address2;
		this.city = json.city;
		this.stateCode = json.state_code;
		this.stateName = json.state_name;
		this.countryCode = json.country_code;
		this.countryName = json.country_name;
		this.postalCode = json.postal_code;
		this.phone = json.phone;
		this.email = json.email;
		this.taxNumber = json.tax_number;
	}

	toJSON() {
		return {
			first_name: this.firstName,
			last_name: this.lastName,
			organization: this.organization,
			address1: this.address1,
			address2: this.address2,
			city: this.city,
			state_code: this.stateCode,
			state_name: this.stateName,
			country_code: this.countryCode,
			country_name: this.countryName,
			postal_code: this.postalCode,
			phone: this.phone,
			email: this.email,
			tax_number: this.taxNumber
		}
	}

	toString() {
		if (this.stateCode) {
			return `${this.address1}, ${this.stateName} ${this.stateCode}, ${this.postalCode}, ${this.countryName}`;
		} else {
			return `${this.address1}, ${this.postalCode}, ${this.countryName}`;
		}
	}

	setFirstName(firstName) {
		this.firstName = firstName;
	}

	getFirstName() {
		return this.firstName ?? '';
	}

	setLastName(lastName) {
		this.lastName = lastName;
	}

	getLastName() {
		return this.lastName ?? '';
	}

	setEmail(email) {
		this.email = email;
	}

	getEmail() {
		return this.email ?? '';
	}

	setAddress1(address1) {
		this.address1 = address1;
	}

	getAddress1() {
		return this.address1 ?? '';
	}

	setAddress2(address2) {
		this.address2 = address2;
	}

	getAddress2() {
		return this.address2 ?? '';
	}

	setCountryCode(countryCode) {
		this.countryCode = countryCode;
	}

	getCountryCode() {
		return this.countryCode ?? '';
	}

	setStateCode(stateCode) {
		this.stateCode = stateCode;
	}

	getStateCode() {
		return this.stateCode ?? '';
	}

	setPostalCode(postalCode) {
		this.postalCode = postalCode;
	}

	getPostalCode() {
		return this.postalCode ?? '';
	}

	setCity(city) {
		this.city = city;
	}

	getCity() {
		return this.city ?? '';
	}
}
