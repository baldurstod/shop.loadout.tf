import { OptionJSON } from '../responses/option';
import { Option, OptionType } from './option';

export class Options {
	#options = new Map<string, Option>();

	add(shopOption: Option): void {
		this.#options.set(shopOption.name, shopOption);
	}

	addOption(name: string, type: OptionType, value: any): void {
		this.add(new Option(name, type, value));
	}

	getOption(optionName: string): Option | null {
		return this.#options.get(optionName) ?? null;
	}

	[Symbol.iterator](): MapIterator<Option> {
		return this.#options.values();
	};

	fromJSON(shopOptionsJson: OptionJSON[] = []): void {
		this.#options.clear();

		for (const shopOptionJson of shopOptionsJson) {
			const shopOption = new Option();
			shopOption.fromJSON(shopOptionJson);
			this.add(shopOption);
		}
	}

	toJSON() {
		const options = [];
		for (const [optionName, option] of this.#options) {
			options.push(option.toJSON());
		}
		return options;
	}
}
