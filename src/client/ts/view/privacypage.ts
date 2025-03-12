import { I18n, createElement, createShadowRoot } from 'harmony-ui';

import commonCSS from '../../css/common.css';
import privacyPageCSS from '../../css/privacypage.css';
import { ShopElement } from './shopelement';

export class PrivacyPage extends ShopElement {
	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [privacyPageCSS, commonCSS],
			child: createElement('div', {
				i18n: {
					innerHTML: '#privacy_policy_content',
				},
			}),
		});
		I18n.observeElement(this.shadowRoot);
	}
}
