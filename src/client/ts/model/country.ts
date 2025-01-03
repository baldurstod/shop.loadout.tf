import { State } from './state';

export class Country {
	#code: string = '';
	#name: string = '';
	#states = new Map<string, State>();
	#region: string = '';

	getCode() {
		return this.#code;
	}

	getName() {
		return this.#name;
	}

	getStates() {
		return this.#states;
	}

	getState(stateCode: string) {
		return this.#states.get(stateCode);
	}

	getRegion() {
		return this.#region;
	}

	hasStates() {
		return this.#states.size > 0;
	}

	fromJSON(countryJSON: any = {}) {
		this.#states.clear();

		this.#code = countryJSON.code;
		this.#name = countryJSON.name;
		this.#region = countryJSON.region;

		const states = countryJSON.states;
		if (states) {
			for (let stateJSON of states) {
				const state = new State();
				state.fromJSON(stateJSON);

				this.#states.set(state.getCode(), state);
			}
		}
	}
}
