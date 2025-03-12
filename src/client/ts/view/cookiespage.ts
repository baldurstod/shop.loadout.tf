import { I18n, createElement, createShadowRoot } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import cookiesPageCSS from '../../css/cookiespage.css';
import { ShopElement } from './shopelement';

export class CookiesPage extends ShopElement {
	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [cookiesPageCSS, commonCSS],
			child: createElement('div', {
				i18n: {
					innerHTML: '#cookies_policy_content',
				},
			}),
		});
		I18n.observeElement(this.shadowRoot);
	}
}
