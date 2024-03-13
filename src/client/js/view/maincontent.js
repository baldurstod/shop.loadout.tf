import { createElement } from 'harmony-ui';

import mainContentCSS from '../../css/maincontent.css';

export class MainContent {
	#htmlElement;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyle: mainContentCSS,
			childs: [
			],
		});
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}
}
