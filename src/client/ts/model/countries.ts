import { CountryJSON } from '../responses/countries';
import { Country } from './country';

export class Countries {
	#countries = new Map<string, Country>();

	[Symbol.iterator]() {
		return this.#countries.values();
	};

	getCountry(countryCode: string) {
		return this.#countries.get(countryCode);
	}

	fromJSON(countriesJSON: Array<CountryJSON>) {
		if (!countriesJSON) {
			return;
		}
		this.#countries.clear();

		for (const countryJSON of countriesJSON) {
			const country = new Country();
			country.fromJSON(countryJSON);
			this.#countries.set(country.getCode(), country);
		}
	}
}
