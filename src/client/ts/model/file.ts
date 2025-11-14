import { FileJSON } from '../responses/product';
;

export class File {
	#type: string;
	#id = 0;
	#url: string;
	//#options: any/*TODO: improve type*/;
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

	/*
	set options(options/*TODO: improve type* /) {
		this.#options = options;
	}

	get options() {
		return this.#options;
	}
	*/

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

	fromJSON(json: FileJSON): void {
		this.type = json.type;
		this.id = json.id;
		this.url = json.url;
		//this.options = json.options as string;
		this.hash = json.hash;
		this.filename = json.filename;
		this.mimeType = json.mime_type;
		this.size = json.size;
		this.width = json.width;
		this.height = json.height;
		this.dpi = json.dpi;
		this.status = json.status;
		this.created = json.created;
		this.thumbnailUrl = json.thumbnail_url;
		this.previewUrl = json.preview_url;
		this.visible = json.visible;
		this.#temporary = json.is_temporary;
	}

	toJSON(): FileJSON {
		return {
			type: this.type,
			id: this.id,
			url: this.url,
			//options: this.options,
			hash: this.hash,
			filename: this.filename,
			mime_type: this.mimeType,
			size: this.size,
			width: this.width,
			height: this.height,
			dpi: this.dpi,
			status: this.status,
			created: this.created,
			thumbnail_url: this.thumbnailUrl,
			preview_url: this.previewUrl,
			visible: this.visible,
			is_temporary: this.#temporary,
		}
	}
}
