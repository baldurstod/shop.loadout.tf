import { State } from './state.js';

export class Country {
	#code;
	#name;
	#states = new Map();
	#region;

	getCode() {
		return this.#code;
	}

	getName() {
		return this.#name;
	}

	getStates() {
		return this.#states;
	}

	getState(stateCode) {
		return this.#states.get(stateCode);
	}

	getRegion() {
		return this.#region;
	}

	hasStates() {
		return this.#states.size > 0;
	}

	fromJSON(countryJSON = {}) {
		this.#states.clear();

		this.#code = countryJSON.code;
		this.#name = countryJSON.name;
		this.#region = countryJSON.region;

		const states = countryJSON.states;
		if (states) {
			for(let stateJSON of states) {
				const state = new State();
				state.fromJSON(stateJSON);

				this.#states.set(state.getCode(), state);
			}
		}
	}
}
