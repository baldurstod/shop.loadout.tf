export enum OptionType {
	None = null,
	Color = 'color',
	Size = 'size',
}

export class Option {
	#name: string;
	#type: OptionType;
	#value: any;

	constructor(name = '', type = OptionType.None, value = '') {
		this.name = name;
		this.type = type;
		this.value = value;
	}

	set name(name) {
		this.#name = name;
	}

	get name() {
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

	set value(value: any) {
		this.#value = value;
	}

	get value(): any {
		return this.#value;
	}

	fromJSON(shopOptionJson: any = {}) {
		this.name = shopOptionJson.name;
		this.type = shopOptionJson.type;
		this.value = shopOptionJson.value;
	}

	toJSON() {
		return {
			name: this.name,
			type: this.type,
			value: this.value,
		};
	}
}
