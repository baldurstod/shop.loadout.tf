import { Country } from "./country";

export class Countries {
	#countries = new Map();

	[Symbol.iterator]() {
		return this.#countries.values();
	};

	getCountry(countryCode) {
		return this.#countries.get(countryCode);
	}

	fromJSON(countriesJSON = []) {
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
