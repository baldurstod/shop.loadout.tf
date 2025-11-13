import { JSONObject } from 'harmony-types';

export class State {
	#code = '';
	#name = '';

	getCode(): string {
		return this.#code;
	}

	getName(): string {
		return this.#name;
	}

	fromJSON(stateJSON: JSONObject): void {
		this.#code = stateJSON.code as string;
		this.#name = stateJSON.name as string;
	}
}
