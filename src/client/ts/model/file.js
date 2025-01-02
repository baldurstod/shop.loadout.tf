export class File  {
	#type;
	#id;
	#url;
	#options;
	#hash;
	#filename;
	#mimeType;
	#size;
	#width;
	#height;
	#dpi;
	#status;
	#created;
	#thumbnailUrl;
	#previewUrl;
	#visible;

	constructor(type, url) {
		this.#type = type;
		this.#url = url;
	}

	set type(type) {
		this.#type = type;
	}

	get type() {
		return this.#type;
	}

	set id(id) {
		this.#id = id;
	}

	get id() {
		return this.#id;
	}

	set url(url) {
		this.#url = url;
	}

	get url() {
		return this.#url;
	}

	set options(options) {
		this.#options = options;
	}

	get options() {
		return this.#options;
	}

	set hash(hash) {
		this.#hash = hash;
	}

	get hash() {
		return this.#hash;
	}

	set filename(filename) {
		this.#filename = filename;
	}

	get filename() {
		return this.#filename;
	}

	set mimeType(mimeType) {
		this.#mimeType = mimeType;
	}

	get mimeType() {
		return this.#mimeType;
	}

	set size(size) {
		this.#size = size;
	}

	get size() {
		return this.#size;
	}

	set width(width) {
		this.#width = width;
	}

	get width() {
		return this.#width;
	}

	set height(height) {
		this.#height = height;
	}

	get height() {
		return this.#height;
	}

	set dpi(dpi) {
		this.#dpi = dpi;
	}

	get dpi() {
		return this.#dpi;
	}

	set status(status) {
		this.#status = status;
	}

	get status() {
		return this.#status;
	}

	set created(created) {
		this.#created = created;
	}

	get created() {
		return this.#created;
	}

	set thumbnailUrl(thumbnailUrl) {
		this.#thumbnailUrl = thumbnailUrl;
	}

	get thumbnailUrl() {
		return this.#thumbnailUrl;
	}

	set previewUrl(previewUrl) {
		this.#previewUrl = previewUrl;
	}

	get previewUrl() {
		return this.#previewUrl;
	}

	set visible(visible) {
		this.#visible = visible;
	}

	get visible() {
		return this.#visible;
	}

	fromJSON(json) {
		this.type = json.type;
		this.id = json.id;
		this.url = json.url;
		this.options = json.options;
		this.hash = json.hash;
		this.filename = json.filename;
		this.mimeType = json.mimeType;
		this.size = json.size;
		this.width = json.width;
		this.height = json.height;
		this.dpi = json.dpi;
		this.status = json.status;
		this.created = json.created;
		this.thumbnailUrl = json.thumbnailUrl;
		this.previewUrl = json.previewUrl;
		this.visible = json.visible;
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
