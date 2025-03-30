import { JSONObject } from '../types';

export class Address {
	firstName = '';
	lastName = '';
	organization = '';
	address1 = '';
	address2 = '';
	city = '';
	stateCode = '';
	stateName = '';
	countryCode = '';
	countryName = '';
	postalCode = '';
	phone = '';
	email = '';
	taxNumber = '';

	get name() {
		return `${this.firstName} ${this.lastName}`;
	}

	fromJSON(json: JSONObject) {
		if (!json) {
			return;
		}
		this.firstName = json.first_name as string;
		this.lastName = json.last_name as string;
		this.organization = json.organization as string;
		this.address1 = json.address1 as string;
		this.address2 = json.address2 as string;
		this.city = json.city as string;
		this.stateCode = json.state_code as string;
		this.stateName = json.state_name as string;
		this.countryCode = json.country_code as string;
		this.countryName = json.country_name as string;
		this.postalCode = json.postal_code as string;
		this.phone = json.phone as string;
		this.email = json.email as string;
		this.taxNumber = json.tax_number as string;
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

	setFirstName(firstName: string) {
		this.firstName = firstName;
	}

	getFirstName() {
		return this.firstName ?? '';
	}

	setLastName(lastName: string) {
		this.lastName = lastName;
	}

	getLastName() {
		return this.lastName ?? '';
	}

	setPhone(phone: string) {
		this.phone = phone;
	}

	getPhone() {
		return this.phone ?? '';
	}

	setEmail(email: string) {
		this.email = email;
	}

	getEmail() {
		return this.email ?? '';
	}

	setAddress1(address1: string) {
		this.address1 = address1;
	}

	getAddress1() {
		return this.address1 ?? '';
	}

	setAddress2(address2: string) {
		this.address2 = address2;
	}

	getAddress2() {
		return this.address2 ?? '';
	}

	setCountryCode(countryCode: string) {
		this.countryCode = countryCode;
	}

	getCountryCode() {
		return this.countryCode ?? '';
	}

	setStateCode(stateCode: string) {
		this.stateCode = stateCode;
	}

	getStateCode() {
		return this.stateCode ?? '';
	}

	setPostalCode(postalCode: string) {
		this.postalCode = postalCode;
	}

	getPostalCode() {
		return this.postalCode ?? '';
	}

	setCity(city: string) {
		this.city = city;
	}

	getCity() {
		return this.city ?? '';
	}
}
