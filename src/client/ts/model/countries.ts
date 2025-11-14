import { CountryJSON } from '../responses/countries';
import { Country } from './country';

export class Countries {
	#countries = new Map<string, Country>();

	[Symbol.iterator](): MapIterator<Country> {
		return this.#countries.values();
	};

	getCountry(countryCode: string): Country | null {
		return this.#countries.get(countryCode) ?? null;
	}

	fromJSON(countriesJSON: CountryJSON[]): void {
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
