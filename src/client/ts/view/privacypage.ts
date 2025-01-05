import { I18n, createElement, createShadowRoot } from 'harmony-ui';

import commonCSS from '../../css/common.css';
import privacyPageCSS from '../../css/privacypage.css';

export class PrivacyPage {
	#shadowRoot?: ShadowRoot;

	#initHTML() {
		this.#shadowRoot = createShadowRoot('section', {
			adoptStyles: [ privacyPageCSS, commonCSS ],
			child: createElement('div', {
				i18n: '#privacy_policy_content',
			}),
		});
		I18n.observeElement(this.#shadowRoot);
		return this.#shadowRoot.host;
	}

	getHTML() {
		return (this.#shadowRoot?.host ?? this.#initHTML()) as HTMLElement;
	}
}
