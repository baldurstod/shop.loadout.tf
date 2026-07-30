import { AddressJSON } from '../responses/order';
import { Address } from './address';

export type UserJSON = {
	display_name: string,
	email: string,
	email_verified: boolean,
	currency: string,
	address: AddressJSON,
}

export class User {
	#displayName = '';
	#email = '';
	#emailVerified = false;
	#currency = 'USD';
	#address = new Address;

	getDisplayName(): string {
		return this.#displayName;
	}

	getEmail(): string {
		return this.#email;
	}

	setEmail(email: string): void {
		this.#email = email;
	}

	isEmailVerified(): boolean {
		return this.#emailVerified;
	}

	setEmailVerified(emailVerified: boolean): void {
		this.#emailVerified = emailVerified;
	}

	getCurrency(): string {
		return this.#currency;
	}

	getAddress(): Address {
		return this.#address.clone();
	}

	fromJSON(json: UserJSON): User {
		this.#displayName = json.display_name;
		this.#email = json.email;
		this.#emailVerified = json.email_verified;
		this.#currency = json.currency;
		this.#address.fromJSON(json.address);
		return this;
	}

	toJSON(): UserJSON {
		return {
			display_name: this.#displayName,
			email: this.#email,
			email_verified: this.#emailVerified,
			currency: this.#currency,
			address: this.#address.toJSON(),
		};
	}

	clone(): User {
		return new User().fromJSON(this.toJSON());
	}
}
