import { State } from './state';

export class Country {
	#code = '';
	#name = '';
	#states = new Map<string, State>();
	#region = '';

	getCode(): string {
		return this.#code;
	}

	getName(): string {
		return this.#name;
	}

	getStates(): Map<string, State> {
		return this.#states;
	}

	getState(stateCode: string): State | null {
		return this.#states.get(stateCode) ?? null;
	}

	getRegion(): string {
		return this.#region;
	}

	hasStates(): boolean {
		return this.#states.size > 0;
	}

	fromJSON(countryJSON: any = {}) {
		this.#states.clear();

		this.#code = countryJSON.code;
		this.#name = countryJSON.name;
		this.#region = countryJSON.region;

		const states = countryJSON.states;
		if (states) {
			for (const stateJSON of states) {
				const state = new State();
				state.fromJSON(stateJSON);

				this.#states.set(state.getCode(), state);
			}
		}
	}
}
