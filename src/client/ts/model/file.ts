import { JSONObject } from 'harmony-types';
;

export class File {
	#type: string;
	#id = 0;
	#url: string;
	#options: any/*TODO: improve type*/;
	#hash = '';
	#filename = '';
	#mimeType = '';
	#size = 0;
	#width = 0;
	#height = 0;
	#dpi = 0;
	#status = '';
	#created = 0;
	#thumbnailUrl = '';
	#previewUrl = '';
	#visible = false;
	#temporary = false;

	constructor(type = '', url = '') {
		this.#type = type;
		this.#url = url;
	}

	set type(type) {
		this.#type = type;
	}

	get type(): string {
		return this.#type;
	}

	set id(id) {
		this.#id = id;
	}

	get id(): number {
		return this.#id;
	}

	set url(url) {
		this.#url = url;
	}

	get url(): string {
		return this.#url;
	}

	set options(options/*TODO: improve type*/) {
		this.#options = options;
	}

	get options() {
		return this.#options;
	}

	set hash(hash) {
		this.#hash = hash;
	}

	get hash(): string {
		return this.#hash;
	}

	set filename(filename) {
		this.#filename = filename;
	}

	get filename(): string {
		return this.#filename;
	}

	set mimeType(mimeType) {
		this.#mimeType = mimeType;
	}

	get mimeType(): string {
		return this.#mimeType;
	}

	set size(size) {
		this.#size = size;
	}

	get size(): number {
		return this.#size;
	}

	set width(width) {
		this.#width = width;
	}

	get width(): number {
		return this.#width;
	}

	set height(height) {
		this.#height = height;
	}

	get height(): number {
		return this.#height;
	}

	set dpi(dpi) {
		this.#dpi = dpi;
	}

	get dpi(): number {
		return this.#dpi;
	}

	set status(status) {
		this.#status = status;
	}

	get status(): string {
		return this.#status;
	}

	set created(created) {
		this.#created = created;
	}

	get created(): number {
		return this.#created;
	}

	set thumbnailUrl(thumbnailUrl) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	get thumbnailUrl(): string {
		return this.#thumbnailUrl;
	}

	set previewUrl(previewUrl) {
		this.#previewUrl = previewUrl;
	}

	get previewUrl(): string {
		return this.#previewUrl;
	}

	set visible(visible) {
		this.#visible = visible;
	}

	get visible(): boolean {
		return this.#visible;
	}

	fromJSON(json: JSONObject): void {
		this.type = json.type as string;
		this.id = json.id as number;
		this.url = json.url as string;
		this.options = json.options as string;
		this.hash = json.hash as string;
		this.filename = json.filename as string;
		this.mimeType = json.mimeType as string;
		this.size = json.size as number;
		this.width = json.width as number;
		this.height = json.height as number;
		this.dpi = json.dpi as number;
		this.status = json.status as string;
		this.created = json.created as number;
		this.thumbnailUrl = json.thumbnailUrl as string;
		this.previewUrl = json.previewUrl as string;
		this.visible = json.visible as boolean;
	}

	toJSON() {
		return {
			type: this.type,
			id: this.id,
			url: this.url,
			options: this.options,
			hash: this.hash,
			filename: this.filename,
			mimeType: this.mimeType,
			size: this.size,
			width: this.width,
			height: this.height,
			dpi: this.dpi,
			status: this.status,
			created: this.created,
			thumbnailUrl: this.thumbnailUrl,
			previewUrl: this.previewUrl,
			visible: this.visible,
		}
	}
}
