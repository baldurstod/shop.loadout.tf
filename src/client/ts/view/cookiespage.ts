import { I18n, createElement, createShadowRoot } from 'harmony-ui';

import commonCSS from '../../css/common.css';
import cookiesPageCSS from '../../css/cookiespage.css';

export class CookiesPage {
	#shadowRoot?: ShadowRoot;

	#initHTML() {
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyles: [ cookiesPageCSS, commonCSS ],
			child: createElement('div', {
				i18n: '#cookies_policy_content',
			}),
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}
}
