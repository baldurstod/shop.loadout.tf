import { File } from './file';

export class Files {
	#files: Array<File> = [];

	get files() {
		return this.#files;
	}

	set files(files) {
		this.#files = files;
	}

	add(file: File) {
		this.#files.push(file);
	}

	addFile(type: string, url: string) {
		this.add(new File(type, url));
	}

	get images() {
		let images = []
		for (let file of this.#files) {
			images.push(file.previewUrl ?? file.url);
		}
		return images;
	}

	getThumbnailUrl(fileType: string) {
		for (let file of this.#files) {
			if (file.type == fileType) {
				return file.thumbnailUrl;
			}
		}
	}

	[Symbol.iterator]() {
		let index = -1;
		const files = this.#files;

		return {
			next: () => ({ value: files[++index], done: !(index in files) })
		};
	};

	fromJSON(filesJSON = []) {
		this.#files = [];

		for (let fileJson of filesJSON) {
			let file = new File();
			file.fromJSON(fileJson);
			this.#files.push(file);
		}
	}

	toJSON() {
		const files = [];
		for (let file of this.#files) {
			files.push(file.toJSON());
		}
		return files;
	}
}
