import { JSONObject } from 'harmony-types';

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

	get name(): string {
		return `${this.firstName} ${this.lastName}`;
	}

	fromJSON(json: JSONObject): void {
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

	toString(): string {
		if (this.stateCode) {
			return `${this.address1}, ${this.stateName} ${this.stateCode}, ${this.postalCode}, ${this.countryName}`;
		} else {
			return `${this.address1}, ${this.postalCode}, ${this.countryName}`;
		}
	}

	setFirstName(firstName: string): void {
		this.firstName = firstName;
	}

	getFirstName(): string {
		return this.firstName ?? '';
	}

	setLastName(lastName: string): void {
		this.lastName = lastName;
	}

	getLastName(): string {
		return this.lastName ?? '';
	}

	setPhone(phone: string): void {
		this.phone = phone;
	}

	getPhone(): string {
		return this.phone ?? '';
	}

	setEmail(email: string): void {
		this.email = email;
	}

	getEmail(): string {
		return this.email ?? '';
	}

	setAddress1(address1: string): void {
		this.address1 = address1;
	}

	getAddress1(): string {
		return this.address1 ?? '';
	}

	setAddress2(address2: string): void {
		this.address2 = address2;
	}

	getAddress2(): string {
		return this.address2 ?? '';
	}

	setCountryCode(countryCode: string): void {
		this.countryCode = countryCode;
	}

	getCountryCode(): string {
		return this.countryCode ?? '';
	}

	setStateCode(stateCode: string): void {
		this.stateCode = stateCode;
	}

	getStateCode(): string {
		return this.stateCode ?? '';
	}

	setPostalCode(postalCode: string): void {
		this.postalCode = postalCode;
	}

	getPostalCode(): string {
		return this.postalCode ?? '';
	}

	setCity(city: string): void {
		this.city = city;
	}

	getCity(): string {
		return this.city ?? '';
	}
}
