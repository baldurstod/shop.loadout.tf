import { I18n, createElement } from 'harmony-ui';

import commonCSS from '../../css/common.css';
import privacyPageCSS from '../../css/privacypage.css';

export class PrivacyPage {
	#htmlElement;

	#initHTML() {
		this.#htmlElement = createElement('section', {
			attachShadow: { mode: 'closed' },
			adoptStyles: [ privacyPageCSS, commonCSS ],
			child: createElement('div', {
				i18n: '#privacy_policy_content',
			}),
		});

		I18n.observeElement(this.#htmlElement);
		return this.#htmlElement;
	}

	get htmlElement() {
		return this.#htmlElement ?? this.#initHTML();
	}
}
