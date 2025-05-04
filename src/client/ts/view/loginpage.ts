import { I18n, createElement, createShadowRoot } from 'harmony-ui';
import commonCSS from '../../css/common.css';
import privacyPageCSS from '../../css/privacypage.css';
import { ShopElement } from './shopelement';
import { Controller } from '../controller';

export class LoginPage extends ShopElement {
	#htmlUsername?: HTMLInputElement;
	#htmlPassword?: HTMLInputElement;

	initHTML() {
		if (this.shadowRoot) {
			return;
		}
		this.shadowRoot = createShadowRoot('section', {
			adoptStyles: [privacyPageCSS, commonCSS],
			childs: [
				createElement('button', {
					innerText: 'login',
					$click: () => Controller.dispatchEvent(new CustomEvent('login', {
						detail: {
							username: this.#htmlUsername?.value,
							password: this.#htmlPassword?.value,
						}
					})),
				}),
				this.#htmlUsername = createElement('input') as HTMLInputElement,
				this.#htmlPassword = createElement('input') as HTMLInputElement,
			]
		});
		I18n.observeElement(this.shadowRoot);
	}
}
