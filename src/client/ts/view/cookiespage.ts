import { I18n, createElement } from 'harmony-ui';

import commonCSS from '../../css/common.css';
import cookiesPageCSS from '../../css/cookiespage.css';

export class CookiesPage {
	#htmlElement;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [ cookiesPageCSS, commonCSS ],
			child: createElement('div', {
				i18n: '#cookies_policy_content',
			}),
		});

		I18n.observeElement(this.#htmlElement);
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}
}
