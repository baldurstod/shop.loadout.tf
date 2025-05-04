import { I18n, createElement, createShadowRoot } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import privacyPageCSS from '../../css/privacypage.css';
import { ShopElement } from './shopelement';
import { Controller } from '../controller';

export class LoginPage extends ShopElement {
	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [privacyPageCSS, commonCSS],
			child: createElement('button', {
				innerText: 'login',
				$click: () => Controller.dispatchEvent(new CustomEvent('login')),
			}),
		});
		I18n.observeElement(this.shadowRoot);
	}
}
