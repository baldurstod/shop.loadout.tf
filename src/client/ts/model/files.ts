import { FileJSON } from '../responses/product';
import { File } from './file';

export class Files {
	#files: File[] = [];

	get files(): File[] {
		return this.#files;
	}

	set files(files) {
		this.#files = files;
	}

	add(file: File): void {
		this.#files.push(file);
	}

	addFile(type: string, url: string): void {
		this.add(new File(type, url));
	}

	get images(): string[] {
		const images = []
		for (const file of this.#files) {
			images.push(file.previewUrl ?? file.url);
		}
		return images;
	}

	getThumbnailUrl(fileType: string): string | null {
		for (const file of this.#files) {
			if (file.type == fileType) {
				return file.thumbnailUrl;
			}
		}
		return null;
	}

	[Symbol.iterator](): Iterator<File> {
		let index = -1;
		const files = this.#files;

		return {
			next: () => ({ value: files[++index]!, done: !(index in files) })
		};
	};

	fromJSON(filesJSON: FileJSON[] = []): void {
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
