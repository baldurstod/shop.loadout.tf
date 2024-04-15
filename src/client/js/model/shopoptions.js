import { ShopOption } from './shopoption.js';

export class ShopOptions {
	#options = new Map();

	constructor() {
	}

	add(shopOption) {
		this.#options.set(shopOption.name, shopOption);
	}

	addOption(name, type, value) {
		this.add(new ShopOption(name, type, value));
	}

	getOption(optionName) {
		return this.#options.get(optionName);
	}

	[Symbol.iterator]() {
		return this.#options.values();
	};

	fromJSON(shopOptionsJson = []) {
		this.#options.clear();

		for (let shopOptionJson of shopOptionsJson) {
			let shopOption = new ShopOption();
			shopOption.fromJSON(shopOptionJson);
			this.add(shopOption);
		}
	}

	toJSON() {
		const options = [];
		for (let [optionName, option] of this.#options) {
			options.push(option.toJSON());
		}
		return options;
	}
}
