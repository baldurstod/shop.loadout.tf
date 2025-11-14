import { OptionJSON } from "../responses/option";

export enum OptionType {
	None = 'none',
	Color = 'color',
	Size = 'size',
}

export class Option {
	#name: string;
	#type: OptionType;
	#value: unknown;

	constructor(name = '', type = OptionType.None, value: unknown = '') {
		this.#name = name;
		this.#type = type;
		this.value = value;
	}

	set name(name) {
		this.#name = name;
	}

	get name(): string {
		return this.#name;
	}

	set type(type: OptionType) {
		/*
		if (type && (type !== 'color' && type !== 'size')) {
			//throw 'Option type must be color or size';
			throw new Error(`Option type must be color or size, got ${type}`);
		}
		*/
		this.#type = type;
	}

	get type(): OptionType {
		return this.#type;
	}

	set value(value: unknown) {
		this.#value = value;
	}

	get value(): unknown {
		return this.#value;
	}

	fromJSON(shopOptionJson: OptionJSON): void {
		this.name = shopOptionJson.name;
		this.type = shopOptionJson.type as OptionType;
		this.value = shopOptionJson.value;
	}

	toJSON(): OptionJSON {
		return {
			name: this.name,
			type: this.type,
			value: String(this.value),
		};
	}
}
