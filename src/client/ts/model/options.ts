import { OptionJSON } from '../responses/option';
import { Option, OptionType } from './option';

export class Options {
	#options = new Map<string, Option>();

	add(shopOption: Option) {
		this.#options.set(shopOption.name, shopOption);
	}

	addOption(name: string, type: OptionType, value: any) {
		this.add(new Option(name, type, value));
	}

	getOption(optionName: string): any {
		return this.#options.get(optionName);
	}

	[Symbol.iterator]() {
		return this.#options.values();
	};

	fromJSON(shopOptionsJson: OptionJSON[] = []) {
		this.#options.clear();

		for (let shopOptionJson of shopOptionsJson) {
			let shopOption = new Option();
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
