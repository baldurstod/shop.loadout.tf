import { FileJSON } from '../responses/product';
import { File } from './file';

export class Files {
	#files: File[] = [];

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
		const images = []
		for (const file of this.#files) {
			images.push(file.previewUrl ?? file.url);
		}
		return images;
	}

	getThumbnailUrl(fileType: string) {
		for (const file of this.#files) {
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

	fromJSON(filesJSON: FileJSON[] = []) {
		this.#files = [];

		for (const fileJson of filesJSON) {
			const file = new File();
			file.fromJSON(fileJson);
			this.#files.push(file);
		}
	}

	toJSON() {
		const files = [];
		for (const file of this.#files) {
			files.push(file.toJSON());
		}
		return files;
	}
}
