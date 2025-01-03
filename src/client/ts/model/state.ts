export class State {
	#code = '';
	#name = '';

	getCode() {
		return this.#code;
	}

	getName() {
		return this.#name;
	}

	fromJSON(stateJSON) {
		this.#code = stateJSON.code;
		this.#name = stateJSON.name;
	}
}
