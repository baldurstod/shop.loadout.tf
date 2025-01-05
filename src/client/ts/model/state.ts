import { JSONObject } from '../types';

export class State {
	#code = '';
	#name = '';

	getCode() {
		return this.#code;
	}

	getName() {
		return this.#name;
	}

	fromJSON(stateJSON: JSONObject) {
		this.#code = stateJSON.code as string;
		this.#name = stateJSON.name as string;
	}
}
